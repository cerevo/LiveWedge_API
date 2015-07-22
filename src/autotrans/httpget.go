package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var (
	content_length string
	last_modified  string
)

func is_changed(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	if resp.Header["Content-Length"] == nil {
		return false
	}
	if resp.Header["Content-Length"][0] != content_length {
		return true
	}
	if resp.Header["Last-Modified"] == nil {
		return false
	}
	if resp.Header["Last-Modified"][0] != last_modified {
		return true
	}
	return false
}

func Get(url string, file string) bool {
	log.Printf("%#v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	ioutil.WriteFile(file, byteArray, 0664)
	if resp.Header["Content-Length"] != nil {
		content_length = resp.Header["Content-Length"][0]
	}
	if resp.Header["Last-Modified"] != nil {
		last_modified = resp.Header["Last-Modified"][0]
	}
	log.Printf("Content-Length: %s\n", content_length)
	log.Printf("Last-Modified: %s\n", last_modified)
	return true
}

func Get_if_changed(url string, file string) bool {
	if is_changed(url) {
		return Get(url, file)
	}
	return false
}

// func main() {
// 	url := "http://172.16.130.126:8081/a.jpg"
// 	file := "_a.jpg"
// 	Get(url, file)
// 	for {
// 		Get_if_changed(url, file)
// 		time.Sleep(2 * time.Second)
// 	}
// }
