package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

func dump(dir string, db *bolt.DB) error {
	filename := path.Join(dir, time.Now().Format("20060102-150405")+".db")
	fmt.Printf("filename=%s\n", filename)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	err = db.View(func(tx *bolt.Tx) error {
		n, err := tx.WriteTo(f)
		log.Printf("DB Dump written %d bytes\n", n)
		return err
	})
	return err
}
