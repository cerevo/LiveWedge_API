package main

import (
	"fmt"
	"libvsw"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])

	c := vsw.RequestSwitcherStatus()
	c2 := vsw.RequestAudioPeakStatus()
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			vsw.HeartBeat()
			log.Printf("status0: HeartBeat!\n")
		case ss := <-c:
			log.Printf("status0: %#v\n", ss)
		case ss := <-c2:
			log.Printf("status0: %#v\n", ss)
		}
	}
}
