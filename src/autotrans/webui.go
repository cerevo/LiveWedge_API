// Copyright 2015, Cerevo Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

const htmlPage string = `
<html><head>
</head><body>
<h1>Auto transition</h1>
<form method="post" action="/">
  <p>Input:<br/>
  <input type="checkbox" name="r" value="1" {{if .Input1}}checked="checked"{{end}} /> 1
  <input type="checkbox" name="r" value="2" {{if .Input2}}checked="checked"{{end}} /> 2
  <input type="checkbox" name="r" value="3" {{if .Input3}}checked="checked"{{end}} /> 3
  <input type="checkbox" name="r" value="4" {{if .Input4}}checked="checked"{{end}} /> 4</p>
  <p>Transition mode:<br/>
  {{select .Trans}} rate {{select .Rate}}</p>
  <p>Interval:<br/>
  {{select .Interval}}</p>
  <input type="checkbox" name="upload" value="true" {{if .UploadStillPicture}}checked="checked"{{end}}/> Upload a still picture and use it as input4<br/>
  File name <input type="text" name="picture" size="40" value="{{.Picture}}" /><br/>
  (Insert a writable SD card to LiveWedge)<br/><br/>
  <input type="checkbox" name="broadcast" value="true" {{if .StartLiveBroadcast}}checked="checked"{{end}} /> Live broadcasting<br/>
  <input type="submit" name="send" value="send" />
  <div align="right"><input type="submit" name="quit" value="quit" /></div>
</form></body></html>`

type tmplParams struct {
	Interval           *forHTMLSelect
	Trans              *forHTMLSelect
	Rate               *forHTMLSelect
	StartLiveBroadcast bool
	UploadStillPicture bool
	Picture string
	Input1             bool
	Input2             bool
	Input3             bool
	Input4             bool
}

type selectItem struct {
	val int
	str string
}

type forHTMLSelect struct {
	Name     string
	Options  []selectItem
	Selected int
}

var (
	params Params
	notify chan Params
	template0 *template.Template
)

var tp = tmplParams{
	Interval: &forHTMLSelect{
		Name: "interval",
		Options: []selectItem{
			selectItem{10, "10 sec"},
			selectItem{30, "30 sec"},
			selectItem{60, "1 min"},
			selectItem{180, "3 min"},
			selectItem{300, "5 min"},
			selectItem{600, "10 min"},
			selectItem{1200, "20 min"},
		},
	},
	Trans: &forHTMLSelect{
		Name: "trans",
		Options: []selectItem{
			selectItem{255, "Wipe (random)"},
			selectItem{0, "Cut"},
			selectItem{1, "Mix"},
			selectItem{2, "Dip to 1"},
			selectItem{258, "Dip to 2"},
			selectItem{514, "Dip to 3"},
			selectItem{770, "Dip to 4"},
			selectItem{3, "Wipe (horizontal)"},
			selectItem{4, "Wipe (horizontal_r)"},
			selectItem{5, "Wipe (vertical)"},
			selectItem{6, "Wipe (vertical_r)"},
			selectItem{7, "Wipe (horizontal_slide)"},
			selectItem{8, "Wipe (horizontal_slide_r)"},
			selectItem{9, "Wipe (vertical_slide)"},
			selectItem{10, "Wipe (vertical_slide_r)"},
			selectItem{11, "Wipe (horizontal_double_slide)"},
			selectItem{12, "Wipe (horizontal_double_slide_r)"},
			selectItem{13, "Wipe (vertical_double_slide)"},
			selectItem{14, "Wipe (vertical_double_slide_r)"},
			selectItem{15, "Wipe (square_top_left)"},
			selectItem{16, "Wipe (square_top_left_r)"},
			selectItem{17, "Wipe (square_top)"},
			selectItem{18, "Wipe (square_top_r)"},
			selectItem{19, "Wipe (square_top_right)"},
			selectItem{20, "Wipe (square_top_right_r)"},
			selectItem{21, "Wipe (square_center_left)"},
			selectItem{22, "Wipe (square_center_left_r)"},
			selectItem{23, "Wipe (square_center)"},
			selectItem{24, "Wipe (square_center_r)"},
			selectItem{25, "Wipe (square_center_right)"},
			selectItem{26, "Wipe (square_center_right_r)"},
			selectItem{27, "Wipe (square_bottom_left)"},
			selectItem{28, "Wipe (square_bottom_left_r)"},
			selectItem{29, "Wipe (square_bottom)"},
			selectItem{30, "Wipe (square_bottom_r)"},
			selectItem{31, "Wipe (square_bottom_right)"},
			selectItem{32, "Wipe (square_bottom_right_r)"},
		},
	},
	Rate: &forHTMLSelect{
		Name: "rate",
		Options: []selectItem{
			selectItem{1000, "1 sec"},
			selectItem{2000, "2 sec"},
			selectItem{3000, "3 sec"},
			selectItem{4000, "4 sec"},
			selectItem{5000, "5 sec"},
			selectItem{6000, "6 sec"},
			selectItem{7000, "7 sec"},
			selectItem{8000, "8 sec"},
			selectItem{9000, "9 sec"},
			selectItem{10000, "10 sec"},
			selectItem{11000, "11 sec"},
			selectItem{12000, "12 sec"},
			selectItem{20000, "20 sec"},
		},
	},
}

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
	params.Picture = r.FormValue("picture")
	notify <- params
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		form(r)
	}

	tp.Interval.Selected = params.Interval
	tp.Trans.Selected = params.Trans
	tp.Rate.Selected = params.Rate
	tp.StartLiveBroadcast = params.StartLiveBroadcast
	tp.UploadStillPicture = params.UploadStillPicture
	tp.Picture = params.Picture
	tp.Input1 = params.Input[0]
	tp.Input2 = params.Input[1]
	tp.Input3 = params.Input[2]
	tp.Input4 = params.Input[3]

	template0.Execute(w, tp)
}

func WebUI(pa Params, nt chan Params) {
	params = pa
	notify = nt

	funcMap := template.FuncMap{
		"select": func(sel *forHTMLSelect) template.HTML {
			h := fmt.Sprintf(`<select name="%s">`, sel.Name)
			for _, v := range sel.Options {
				var s string
				if v.val == sel.Selected {
					s = " selected"
				}
				h += fmt.Sprintf(`<option value="%d"%s>%s</option>`, v.val, s, v.str)
			}
			h += "</select>"
			return template.HTML(h)
		},
	}
	template0 = template.Must(template.New("autotrans").Funcs(funcMap).Parse(htmlPage))

	fmt.Printf("Open http://localhost:8080/ by web browser.\n")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
