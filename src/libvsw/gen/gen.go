package main

import (
	"fmt"
	"log"
	"os"
)

const head = `// This file is auto-generated by gen.go. Do not edit.

package libvsw

import (
	"bytes"
	"encoding/binary"
	"log"
	"unsafe"
)
`

const fa = `
var (
	%[1]s %[1]sType
	_%[1]sChan chan %[1]sType
)

func (vsw Vsw) Request%[1]s() <-chan %[1]sType {
	if _%[1]sChan == nil {
		_%[1]sChan = make(chan %[1]sType)
	}
	return _%[1]sChan
}

func read%[1]s(len int, reader *bytes.Reader) {
	var a %[1]sType
	if len != int(unsafe.Sizeof(a)) {
		log.Printf(" size mismatch %%T len=%%d\n", a, len)
		return
	}
	err := binary.Read(reader, LE, &a)
	checkError(err)
	if %[1]s != a {
		//log.Printf("%%#v\n", a)
		%[1]s = a
	        if _%[1]sChan != nil {
		       _%[1]sChan <- a
	        }
	}
}
`

func main() {
	a := []string{
		"SwMode",
		"DefaultBackgroundColor",
		"MountStatus",
		"PreviewMode",
		"RecordingStatus",
		"FadeToDefaultColorStatus",
		"ExternalInputStatus",
		"ProgramOutStatus",
		"PreviewOutStatus",
		"CasterMessage",
		"CasterStatistics",
		"SwitcherStatus",
		"AudioMixerStatus",
		"AudioMixerAllStatus",
		"AudioPeakStatus",
	}

	file, err := os.Create("a.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprint(file, head)
	for _, v := range a {
		fmt.Fprintf(file, fa, v)
	}
}
