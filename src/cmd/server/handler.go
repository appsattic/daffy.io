package main

import "net/http"

func serveFile(filename string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

func fileServer(dirname string) http.Handler {
	return http.FileServer(http.Dir(dirname))
}
