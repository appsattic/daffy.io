# mux : powerful router, with first-class middleware chains.

## Overview [![GoDoc](https://godoc.org/github.com/gomiddleware/mux?status.svg)](https://godoc.org/github.com/gomiddleware/mux) [![Build Status](https://travis-ci.org/gomiddleware/mux.svg)](https://travis-ci.org/gomiddleware/mux)

Instead of focussing on pure-speed trie based router implementations, mux instead focusses on being both small yet
powerful. Some of the main features of mux are features that have been left out, such as:

* no router groups
* no sub-router mouting
* no ignoring case on paths
* no automatic slash/non-slash redirection

The features that mux boasts are all idomatic Go, such as:

* uses the standard context package
* middleware for route prefixes
* middleware chains for all router endpoints
* no external dependencies, just plain net/http
* everything is explicit - and is very much considered a feature

The combination of just these two things give you a very powerful composition system where you compose middleware on
prefixes and middleware chains on endpoints.

## Installation

```sh
go get github.com/gomiddleware/mux
```

## Usage / Example

```go
// new Mux
r := mux.New()

// every request gets a 'X-Request-ID' request header
r.Use("/", reqid.RandomId)

// serve the /about page
r.Get("/about", nil, http.HandlerFunc(aboutHandler))

// log requests to the homepage
r.Get("/", []func(http.Handler) http.Handler{requestLogger}, http.HandlerFunc(homeHandler))

// start the server
http.ListenAndServe(":8080", r)
```

## Author ##

By [Andrew Chilton](https://chilts.org/), [@twitter](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com/), [@AppsAttic](https://twitter.com/AppsAttic).

## LICENSE

MIT.
