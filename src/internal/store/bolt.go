package store

import (
	"time"

	"github.com/boltdb/bolt"
)

type BoltStore struct {
	filename string
	db       *bolt.DB
}

func NewBoltStore(filename string) *BoltStore {
	return &BoltStore{
		filename: filename,
	}
}

func (b *BoltStore) Open() error {
	// open the db
	db, err := bolt.Open(b.filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	b.db = db
	return err
}

func (b *BoltStore) Close() error {
	return b.db.Close()
}
