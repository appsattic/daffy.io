package types

import "time"

type Social struct {
	Id                string // e.g. "twitter:123456"
	UserId            string // e.g. "de58631b-fd37-40a4-8573-c96acd7ed22e" - the FK to our Users
	NickName          string // e.g. "andychilton" - the NickName they have from the Social Provider
	Title             string // e.g. "Andrew Chilton" - the title we got from the Social Provider
	Email             string // e.g. "andychilton@gmail.com" - the email we got from the Social Provider
	AccessToken       string // e.g. "deafbeef"
	AccessTokenSecret string // e.g. "cafebabe"
	RefreshToken      string // e.g. "baadf00d"
	Inserted          time.Time
	Updated           time.Time
}
