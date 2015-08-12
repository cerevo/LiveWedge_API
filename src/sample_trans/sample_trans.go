package main

import (
	"fmt"
	"libvsw"
	"log"
	"os"
	"time"
)

func sample_trans(vsw *libvsw.Vsw, src1, src2 int) {
	log.Printf("sample_trans: Cut to input %d.\n", src1)
	vsw.Cut(src1)
	time.Sleep(3 * time.Second)
	rate := 1000
	log.Printf("sample_trans: Mix to input %d at rate %d msec.\n", src2, rate)
	vsw.Mix(rate, src2)
	time.Sleep(3 * time.Second)
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

	for {
		sample_trans(vsw, 1, 4)
	}
}
