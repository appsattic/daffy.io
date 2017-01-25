package main

import (
	"net/http"

	"github.com/gomiddleware/slash"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func main() {
	handle := http.HandlerFunc(handler)

	http.Handle("/", handle)

	http.Handle("/docs", slash.Add(handle))
	http.Handle("/docs/", handle)

	http.Handle("/about", handle)
	http.Handle("/about/", slash.Remove(handle))

	http.ListenAndServe(":8080", nil)
}
