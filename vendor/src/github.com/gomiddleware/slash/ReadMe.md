# GoMiddleware : Slash #

Middleware that redirects to URLs either with or without a trailing slash.

* [Project](https://github.com/gomiddleware/slash)
* [GoDoc](https://godoc.org/github.com/gomiddleware/slash)

## Synopsis ##

```go
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
```

## Author ##

Written by [Andrew Chilton](https://chilts.org/) for [Apps Attic Ltd](https://appsattic.com/).

## License ##

ISC.

(Ends)

