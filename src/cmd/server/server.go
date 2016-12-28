package main

import (
	"fmt"
	"log"

	"internal/store"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkInterface(store store.Api) {
	fmt.Printf("store=%#v\n", store)
}

func main() {
	// create/open/connect to a store
	boltStore := store.NewBoltStore("webapp.db")
	err := boltStore.Open()
	check(err)
	defer boltStore.Close()

	checkInterface(boltStore)

	fmt.Printf("Hello, World!\n")
}
