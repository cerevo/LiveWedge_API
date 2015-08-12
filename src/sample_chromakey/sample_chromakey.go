package main

import (
	"fmt"
	"libvsw"
	//"log"
	"os"
	"time"
)

func sample_chromakey(vsw *libvsw.Vsw) {
	wait := 4000

	vsw.UploadFile("greenback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_GREEN)
	vsw.CutUs(1, 4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.CutUs(1, 0)

	vsw.UploadFile("blueback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_BLUE)
	vsw.CutUs(1, 4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.CutUs(1, 0)

	vsw.UploadFile("purpleback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_PURPLE)
	vsw.CutUs(1, 4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.CutUs(1, 0)

	vsw.UploadFile("redback.jpg")
	vsw.SetChromaKey(libvsw.CHROMA_KEY_RED)
	vsw.CutUs(1, 4)
	time.Sleep(time.Duration(wait) * time.Millisecond)
	vsw.CutUs(1, 0)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])
	defer vsw.Close()
	sample_chromakey(vsw)
}
