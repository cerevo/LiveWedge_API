package main

import (
	"libvsw"
	"fmt"
)

func main() {
	items := libvsw.Find()
	fmt.Printf("Found %d LiveWedge(s).\n", len(items))
	for _, v := range items {
		fmt.Printf("%#v\n", v)
	}
}
