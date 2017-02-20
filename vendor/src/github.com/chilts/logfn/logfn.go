// Package logfn just allows you to do this:
//
//     func myFunc() {
//         defer logfn.Exit(logfn.Enter("myFunc"))
//     }
//
// Which will appear in the logs as:
//
//     2017/02/20 22:49:05 -> main.myFunc
//     2017/02/20 22:49:05 <- main.myFunc
//
// It uses the global log, however it probably should create a type with an enter and exit method which takes a log to
// be used. PRs welcome. :)
package logfn

import "log"

// Enter just logs "-> name" and returns name so it can be used in the defer logfn.Exit(name).
func Enter(name string) string {
	log.Println("-> " + name)
	return name
}

// Exit just logs "<- name".
func Exit(name string) {
	log.Println("<- " + name)
}
