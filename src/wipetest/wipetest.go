package main

import (
	"fmt"
	"libvsw"
	"log"
	"os"
	"time"
)

func wipeTest(vsw *libvsw.Vsw, src1, src2 int) {
	log.Printf("wipeTest src1=%d src2=%d\n", src1, src2)
	vsw.Cut(src1)

	for wipe_type := 0; wipe_type < libvsw.WIPE_TYPE_NUM; wipe_type++ {
		log.Printf("wipe_type=%d\n", wipe_type)
		vsw.Wipe(2000, src2, wipe_type)
		time.Sleep(3 * time.Second)
		vsw.Wipe(2000, src1, wipe_type)
		time.Sleep(3 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw := libvsw.NewVsw(os.Args[1])

	for {
		wipeTest(vsw, 1, 4)
	}
}
