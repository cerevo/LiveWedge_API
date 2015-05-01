package main

import (
	"fmt"
	"libvsw"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])

	c := vsw.RequestSwitcherStatus()
	c2 := vsw.RequestAudioPeakStatus()
	for {
		vsw.HeartBeat()
		select {
		case ss := <-c:
			fmt.Printf("status0: %#v\n", ss)
		case ss := <-c2:
			fmt.Printf("status0: %#v\n", ss)
		}
	}
}
