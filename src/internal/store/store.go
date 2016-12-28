package store

type Api interface {
	Open() error
	Close() error
}
