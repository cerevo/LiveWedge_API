package main

import (
	"fmt"
	"libvsw"
	"log"
	"os"
	"time"
)

func sample_status(vsw *libvsw.Vsw) {
	c := vsw.RequestSwitcherStatus()
	c2 := vsw.RequestRecordingResult()
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			vsw.HeartBeat()
			log.Printf("sample_status: HeartBeat!\n")
		case ss := <-c:
			log.Printf("sample_status: %#v\n", ss)
		case ss := <-c2:
			log.Printf("sample_status: %#v\n", ss)
		}
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

	sample_status(vsw)
}
