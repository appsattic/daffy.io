package types

import "time"

type User struct {
	Id        string   // e.g. "de58631b-fd37-40a4-8573-c96acd7ed22e"
	Name      string   // e.g. "chilts" (unique)
	Title     string   // e.g. "Andrew Chilton"
	Email     string   // e.g. "andychilton@gmail.com"
	SocialIds []string // e.g. [ "twitter:123456", "facebook:123" ]
	Inserted  time.Time
	Updated   time.Time
}

type UpdateUser struct {
	Name  string `schema:"userName" valid:"required,length(3|32),matches(^[a-z][a-z0-9-]+[a-z0-9]$)"`
	Title string `schema:"title" valid:"required"`
	Email string `schema:"email" valid:"required,email"`
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
