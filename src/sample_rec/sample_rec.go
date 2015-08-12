package main

import (
	"fmt"
	"libvsw"
	"net/http"
	"os"
)

const f0 string = `<html><body>
Recording
<form method="post" action="/">
  <input type="submit" name="rec" value="start" />
  <input type="submit" name="rec" value="stop" />
</form>
Just send command without any error checking.
</body></html>`

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
	fmt.Printf("Open http://localhost:8080/ by web browser.\n")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("rec") == "start" {
			vsw.ChangeRecordingState(1)
		} else if r.FormValue("rec") == "stop" {
			vsw.ChangeRecordingState(0)
		}
		fmt.Fprint(w, f0)
	})
	http.ListenAndServe(":8080", nil)
}
