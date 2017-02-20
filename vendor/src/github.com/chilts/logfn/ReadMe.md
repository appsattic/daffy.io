# logfn : simple helper with logging function enter/exit

## Overview [![GoDoc](https://godoc.org/github.com/chilts/logfn?status.svg)](https://godoc.org/github.com/chilts/logfn) [![Build Status](https://travis-ci.org/chilts/logfn.svg?branch=master)](https://travis-ci.org/chilts/logfn) [![Code Climate](https://codeclimate.com/github/chilts/logfn/badges/gpa.svg)](https://codeclimate.com/github/chilts/logfn) [![Go Report Card](https://goreportcard.com/badge/github.com/chilts/logfn)](https://goreportcard.com/report/github.com/chilts/logfn) [![Sourcegraph](https://sourcegraph.com/github.com/chilts/logfn/-/badge.svg)](https://sourcegraph.com/github.com/chilts/logfn?badge)

Simple helper which allows you to easily log function enter/exit lines.

## Install

```
go get github.com/chilts/logfn
```

## Example

```
package main

import "github.com/chilts/logfn"

func main() {
	defer logfn.Exit(logfn.Enter("main"))
	myFunc()
}

func myFunc() {
	defer logfn.Exit(logfn.Enter("myFunc"))
}
```

Will output lines like:

```
2017/02/20 23:24:44 -> main
2017/02/20 23:24:44 -> myFunc
2017/02/20 23:24:45 <- myFunc
2017/02/20 23:24:45 <- main
```

## Author

By [Andrew Chilton](https://chilts.org/), [@twitter](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com/), [@AppsAttic](https://twitter.com/AppsAttic).

## License

[MIT](https://publish.li/mit-qLQqmVTO).

(Ends)
