// Package main just shows an example program. It has two functions (including main) and should just log enter/exit for
// both. Example output is:
//
//     2017/02/20 23:24:44 -> main
//     2017/02/20 23:24:44 -> myFunc
//     2017/02/20 23:24:44 <- myFunc
//     2017/02/20 23:24:44 <- main
//
// To run it, just try `go run example.go`.
package main

import "github.com/chilts/logfn"

func main() {
	defer logfn.Exit(logfn.Enter("main"))
	myFunc()
}

func myFunc() {
	defer logfn.Exit(logfn.Enter("myFunc"))
}
