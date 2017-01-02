package types

import "time"

type Social struct {
	Id       string // e.g. "twitter:123456"
	UserName string // e.g. "chilts" - the FK to our internal users
	NickName string // e.g. "andychilton" - the NickName they have from the Social Provider
	Title    string // e.g. "Andrew Chilton" - the title we got from the Social Provider
	Email    string // e.g. "andychilton@gmail.com" - the email we got from the Social Provider
	Inserted time.Time
	Updated  time.Time
}
