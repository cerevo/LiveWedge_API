package main

import (
	"fmt"
	"libvsw"
	//"log"
	"os"
	"time"
)

func pinpTest(vsw *libvsw.Vsw, mainSrc, subSrc int) {
	vsw.PinpGeometry(0, 0, 65536/4, 65536/4, 0, 0, 65536, 65536)
	vsw.PinpBorder(0x00ff00ff, 4)
	vsw.SetSubMode(libvsw.SUB_MODE_PINP)
	vsw.CutUs(mainSrc, subSrc)
	rate := 1000
	for {
		for i := 1; i < 4; i++ {
			time.Sleep(time.Duration(rate) * time.Millisecond)
			vsw.PinpGeometry(65536/4*i, 0, 65536/4, 65536/4, 0, 0, 65536, 65536)
		}
		for i := 1; i < 4; i++ {
			time.Sleep(time.Duration(rate) * time.Millisecond)
			vsw.PinpGeometry(65536/4*3, 65536/4*i, 65536/4, 65536/4, 0, 0, 65536, 65536)
		}
		for i := 2; i >= 0; i-- {
			time.Sleep(time.Duration(rate) * time.Millisecond)
			vsw.PinpGeometry(65536/4*i, 65536/4*3, 65536/4, 65536/4, 0, 0, 65536, 65536)
		}
		for i := 2; i >= 0; i-- {
			time.Sleep(time.Duration(rate) * time.Millisecond)
			vsw.PinpGeometry(65536/4*0, 65536/4*i, 65536/4, 65536/4, 0, 0, 65536, 65536)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])

	vsw.Cut(1)
	pinpTest(vsw, 1, 4)
}
