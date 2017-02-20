package main

import "log"

func enter(name string) string {
	log.Println("-> " + name)
	return name
}

func exit(name string) {
	log.Println("<- " + name)
}
