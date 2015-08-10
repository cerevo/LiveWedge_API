package main

import (
	"encoding/json"
	"fmt"
	"libvsw"
	"log"
	"os"
	"time"
	"math/rand"
)

const PARAM_VERSION = 5
const PARAMS_FILE = "autotrans.json"
const TMP_PICTURE = "_a.jpg"
const DEFAULT_PICTURE = "a.jpg"

type Params struct {
	Param_version      int
	Input              [4]bool
	Interval           int /* sec */
	Trans              int /* bit0-7 {CUT, MIX, DIP, WIPE, ..}, bit8-9:dip_src */
	Rate               int /* msec */
	StartLiveBroadcast bool
	UploadStillPicture bool
	PictureUrl string
}

var defaultParams = Params{
	Param_version:      PARAM_VERSION,
	Input:              [4]bool{true, true, false, false},
	Interval:           30,
	Trans:              0x302, /* Dip to 4 */
	Rate:               5000,
	StartLiveBroadcast: false,
	UploadStillPicture: false,
	PictureUrl: "",
}

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func loop(vsw *libvsw.Vsw, pa Params, notify chan Params) {
	index := 0
	for {
		select {
		case pa = <-notify:
			log.Printf("got from chan\n")
			saveParams(pa)
		case <-time.After(time.Second * time.Duration(pa.Interval)):
			log.Printf("periodic timer\n")
		}
		if pa.UploadStillPicture {
			if Get_if_changed(pa.PictureUrl, TMP_PICTURE) {
				vsw.UploadFile(TMP_PICTURE)
			}
		}
		index = (index + 1) % 4
		i := 0
		for ; i < 4; i++ {
			if pa.Input[i] == true {
				break
			}
		}
		if i == 4 {
			// no input checked
			pa.Interval = 1000000
			continue
		}
		for pa.Input[index] == false {
			index = (index + 1) % 4
		}
		switch pa.Trans & 0xff {
		case 0:
			vsw.Cut(index + 1)
		case 1:
			vsw.Mix(pa.Rate, index+1)
		case 2:
			vsw.Dip(pa.Rate, index+1, ((pa.Trans>>8)&3)+1)
		case 255:
			vsw.Wipe(pa.Rate, index+1, random(0, libvsw.WIPE_TYPE_NUM - 1))
		default:
			vsw.Wipe(pa.Rate, index+1, pa.Trans-libvsw.TRANSITION_TYPE_WIPE)
		}
	}
}

func loadParams(filename string) Params {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return defaultParams
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var pa Params
	err = dec.Decode(&pa)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return defaultParams
	}
	if pa.Param_version != PARAM_VERSION {
		fmt.Fprintf(os.Stderr, "Param_version mismatch. Use default.\n")
		return defaultParams
	}
	return pa
}

func saveParams(pa Params) {
	f, err := os.OpenFile(PARAMS_FILE, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(pa)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])
	vsw.Cut(1)

	pa := loadParams(PARAMS_FILE)
	if pa.StartLiveBroadcast {
		vsw.ChangeLiveBroadcastState(1)
	}
	if pa.UploadStillPicture {
		if Get(pa.PictureUrl, TMP_PICTURE) {
			vsw.UploadFile(TMP_PICTURE)
		} else {
			vsw.UploadFile(DEFAULT_PICTURE)
		}
	}

	notify := make(chan Params, 1)
	go loop(vsw, pa, notify)

	WebUI(pa, notify)
}
