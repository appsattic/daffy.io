package types

import "time"

type User struct {
	Name      string   // e.g. "chilts" (ie. their Twitter handle)
	Title     string   // e.g. "Andrew Chilton"
	Email     string   // e.g. "andychilton@gmail.com"
	SocialIds []string // e.g. [ "twitter-123456", "facebook-123" ]
	Inserted  time.Time
	Updated   time.Time
}

// Validate firstly normalises the thing, then validates it and returns either true (valid) or false (invalid). It sets any messages onto
// the Thing.Error field for display.
func (x *User) Validate() bool {
	now := time.Now().UTC()

	// normalise
	x.Inserted = now
	x.Updated = now

	// ToDo: check Name contains only 'a-z0-9-', min 3 chars, max 32, starts with a letter.

	return true
}
