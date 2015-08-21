// Copyright 2015, Cerevo Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"libvsw"
	"log"
	"math/rand"
	"os"
	"time"
)

const PARAM_VERSION = 5
const PARAMS_FILE = "autotrans.json"

type Params struct {
	Param_version      int
	Input              [4]bool
	Interval           int /* sec */
	Trans              int /* bit0-7 {CUT, MIX, DIP, WIPE, ..}, bit8-9:dip_src */
	Rate               int /* msec */
	StartLiveBroadcast bool
	UploadStillPicture bool
	Picture            string
}

var defaultParams = Params{
	Param_version:      PARAM_VERSION,
	Input:              [4]bool{true, true, false, false},
	Interval:           30,
	Trans:              255, /* Random wipe */
	Rate:               3000,
	StartLiveBroadcast: false,
	UploadStillPicture: false,
	Picture:            "",
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func loop(vsw *libvsw.Vsw, pa Params, notify chan Params) {
	index := 0
	for {
		select {
		case pa = <-notify:
			log.Printf("got from chan\n")
			saveParams(pa)
			if pa.StartLiveBroadcast {
				vsw.ChangeLiveBroadcastState(1)
			} else {
				vsw.ChangeLiveBroadcastState(0)
			}
			if pa.UploadStillPicture {
				vsw.UploadFile(pa.Picture)
			}
		case <-time.After(time.Second * time.Duration(pa.Interval)):
			log.Printf("periodic timer\n")
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
			vsw.Wipe(pa.Rate, index+1, random(0, libvsw.WIPE_TYPE_NUM-1))
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
		fmt.Fprintf(os.Stderr, "usage: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw, err := libvsw.NewVsw(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open LiveWedge: %s\n", err)
		os.Exit(1)
	}
	defer vsw.Close()

	vsw.Cut(1)

	pa := loadParams(PARAMS_FILE)
	if pa.StartLiveBroadcast {
		vsw.ChangeLiveBroadcastState(1)
	} else {
		vsw.ChangeLiveBroadcastState(0)
	}
	if pa.UploadStillPicture {
		vsw.UploadFile(pa.Picture)
	}

	notify := make(chan Params, 1)
	go loop(vsw, pa, notify)

	WebUI(pa, notify)
}
