package main

import (
	"fmt"
	"libvsw"
	"net/http"
	"os"
)

var vsw *libvsw.Vsw

const f0 string = `<html><body>
Recording
<form method="post" action="/">
  <input type="submit" name="rec" value="start" />
  <input type="submit" name="rec" value="stop" />
</form>
Just send command without any error checking.
</body></html>`

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("rec=%s\n", r.FormValue("rec"))

	if r.FormValue("rec") == "start" {
		vsw.ChangeRecordingState(1)
	} else if r.FormValue("rec") == "stop" {
		vsw.ChangeRecordingState(0)
	}
	fmt.Fprint(w, f0)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usge: %s IP_address_of_livewedge\n", os.Args[0])
		os.Exit(1)
	}
	vsw = libvsw.NewVsw(os.Args[1])

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
