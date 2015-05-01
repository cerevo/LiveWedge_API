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

	c := libvsw.RequestSwitcherStatus()
	c2 := libvsw.RequestAudioPeakStatus()
	for {
		vsw.HeartBeat()
		select {
		case ss := <-c:
			fmt.Printf("Got from chan! %v", ss)
		case ss := <-c2:
			fmt.Printf("Got from chan! %v", ss)
		}
	}
}
