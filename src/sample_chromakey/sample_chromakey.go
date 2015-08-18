package main

import (
	"fmt"
	"libvsw"
	//"log"
	"os"
	"time"
)

func sample_chromakey(vsw *libvsw.Vsw) {
	wait := 8000
	vsw.CutSub(0)
	vsw.Cut(1)

	// You have to insert a writable SD card to LiveWedge.
	vsw.UploadFile("greenback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_GREEN)
	vsw.CutSub(4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.CutSub(0)

	vsw.UploadFile("blueback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_BLUE)
	vsw.MixSub(2000, 4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.MixSub(2000, 0)

	vsw.UploadFile("purpleback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_PURPLE)
	vsw.CutSub(4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.CutSub(0)

	vsw.UploadFile("redback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_RED)
	vsw.MixSub(2000, 4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.MixSub(2000, 0)
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
	sample_chromakey(vsw)
}
