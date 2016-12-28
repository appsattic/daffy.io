package types

import "time"

type User struct {
	Name     string // e.g. "chilts" (ie. their Twitter handle)
	Title    string // e.g. "Andrew Chilton"
	Email    string // e.g. "andychilton@gmail.com"
	Inserted time.Time
	Updated  time.Time
}
