package main

import (
	"io/ioutil"
	"net/http"
)

func HandleFavicon(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("static/favicon.ico")
	if err != nil {
		panic(err)
	}
	w.Write(dat)
}
func HandleStaticAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Handled-Path", "."+r.URL.Path)
	w.Header().Set("Server", "loog")
	http.ServeFile(w, r, "."+r.URL.Path)
}
