package main

import (
	"fmt"
	"net/http"
	"strconv"
	"os"
)

const f0 string = `<html><head>
</head><body>
<h1>Auto transition</h1>
<h3>Current Setting:</h3>
`
const f1 string = `
<hr/>
<form method="post" action="/">
  <p>Input:<br/>
  <input type="checkbox" name="r" value="1" /> 1
  <input type="checkbox" name="r" value="2" /> 2
  <input type="checkbox" name="r" value="3" /> 3
  <input type="checkbox" name="r" value="4" /> 4</p>
  <p>Transition mode:<br/>
  <select name="trans">
    <option value="0"></option>
    <option value="0">Cut</option>
    <option value="1">Mix</option>
    <option value="2">Dip to 1</option>
    <option value="18">Dip to 2</option>
    <option value="34">Dip to 3</option>
    <option value="50">Dip to 4</option>
  </select> rate <select name="rate">
    <option value="5000"></option>
    <option value="1000">1 sec</option>
    <option value="2000">2 sec</option>
    <option value="3000">3 sec</option>
    <option value="4000">4 sec</option>
    <option value="5000">5 sec</option>
    <option value="6000">6 sec</option>
    <option value="10000">10 sec</option>
  </select></p>
  <p>Interval:<br/>
  <select name="interval">
    <option value="30"></option>
    <option value="10">10 sec</option>
    <option value="30">30 sec</option>
    <option value="60">1 min</option>
    <option value="180">3 min</option>
    <option value="600">10 min</option>
    <option value="1200">20 min</option>
  </select></p>
  <p>Boot time options:<br/>
  <input type="checkbox" name="broadcast" value="true" /> Start live broadcasting<br/>
 <input type="checkbox" name="upload" value="true" /> Upload a still picture and use it as input4<br/>
  <input type="submit" name="send" value="send" />
  <div align="right"><input type="submit" name="quit" value="quit" /></div>
</form>
</body></html>`

var (
	params Params
	notify chan Params
)

func form(r *http.Request) {
	if r.FormValue("quit") == "quit" {
		os.Exit(0)
	}
	if r.FormValue("send") != "send" {
		return
	}
	r.ParseForm()
	for i := range params.Input {
		params.Input[i] = false
	}
	for _, i := range r.Form["r"] {
		val, err := strconv.ParseInt(i, 10, 32)
		if err == nil {
			params.Input[val-1] = true
		}
	}
	interval, _ := strconv.ParseInt(r.FormValue("interval"), 10, 32)
	t, _ := strconv.ParseInt(r.FormValue("trans"), 10, 32)
	rate, _ := strconv.ParseInt(r.FormValue("rate"), 10, 32)
	params.Interval = int(interval)
	params.Trans = int(t)
	params.Rate = int(rate)
	params.StartLiveBroadcast = (r.FormValue("broadcast") == "true")
	params.UploadStillPicture = (r.FormValue("upload") == "true")
	notify <- params

}
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form(r)
	}
	fmt.Fprint(w, f0)
	for i, v := range params.Input {
		if v {
			fmt.Fprintf(w, "input[%d] = ON</br>", i+1)
		} else {
			fmt.Fprintf(w, "input[%d] = OFF</br>", i+1)
		}
	}
	fmt.Fprintf(w, "transition mode: ")
	switch params.Trans & 3 {
	case 0:
	   fmt.Fprintf(w, "cut<br/>")
	case 1:
	   fmt.Fprintf(w, "mix, rate: %d msec<br/>", params.Rate)
	case 2:   
	   fmt.Fprintf(w, "dip to %d, rate: %d msec<br/>", ((params.Trans >> 4) & 3) + 1, params.Rate)
        }
	fmt.Fprintf(w, "interval: %d sec<br/>", params.Interval)
	fmt.Fprintf(w, "<br/>Boot time options:<br/>")
	fmt.Fprintf(w, "  Start live broadcasting: %t<br/>", params.StartLiveBroadcast)
	fmt.Fprintf(w, "  Upload a still picture and use it as input4: %t<br/>", params.UploadStillPicture)

	fmt.Fprint(w, f1)
}

func WebUI(pa Params, nt chan Params) {
	params = pa
	notify = nt
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
