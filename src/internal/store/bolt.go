package store

import (
	"fmt"
	"time"

	"github.com/Machiel/slugify"
	"github.com/boltdb/bolt"
	"github.com/chilts/rod"
	uuid "github.com/hashicorp/go-uuid"

	"internal/types"
)

var userBucket = "user"
var socialBucket = "social"
var indexUserNameUniqueIndex = "i-u-n-u"

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

func (b *BoltStore) LogIn(provider, id, nickName, title, email string) (*types.User, error) {
	var user types.User
	var social types.Social
	now := time.Now().UTC()

	// 1. see if this social id exists
	// 2. if it does, read the user and return it
	// 3. if it doesn't, add the Social and User types

	fmt.Printf("boltStore.LogIn(): entry\n")
	fmt.Printf("* provider=%#v\n", provider)
	fmt.Printf("* id=%#v\n", id)
	fmt.Printf("* nickName=%#v\n", nickName)
	fmt.Printf("* title=%#v\n", title)
	fmt.Printf("* email=%#v\n", email)

	err := b.db.Update(func(tx *bolt.Tx) error {
		// create a socialId that we use internally (to look the user up)
		socialId := provider + ":" + id

		// fetch this Social entity
		errGetSocial := rod.GetJson(tx, socialBucket, socialId, &social)
		if errGetSocial != nil {
			return errGetSocial
		}

		// check to see if the socialId exists
		if social.Id != "" {
			fmt.Printf("Got social = %#v\n", social)
			// get this user - should ALWAYS work if the above Social exists
			errGetUser := rod.GetJson(tx, userBucket, social.UserId, &user)
			fmt.Printf("Got user = %#v\n", user)
			if errGetUser != nil {
				return errGetUser
			}
			return nil
		}

		// create a unique UserId for this user (this never changes)
		userId, _ := uuid.GenerateUUID()

		// create the Social
		social = types.Social{
			Id:       socialId,
			UserId:   userId,
			NickName: nickName,
			Title:    title,
			Email:    email,
			Inserted: now,
			Updated:  now,
		}
		fmt.Printf("Adding a new Social = %#v\n", social)
		errPutSocial := rod.PutJson(tx, socialBucket, socialId, social)
		if errPutSocial != nil {
			return errPutSocial
		}

		// create a unique userName for this user - they can change it if they like
		userName := slugify.Slugify(nickName + "-" + id)

		// create the User
		user = types.User{
			Id:    userId,
			Name:  userName,
			Title: title,
			Email: email,
			SocialIds: []string{
				socialId,
			},
			Inserted: now,
			Updated:  now,
		}
		fmt.Printf("Adding a new User = %#v\n", user)
		errPutUser := rod.PutJson(tx, userName, userId, user)
		if errPutUser != nil {
			return errPutUser
		}

		// index the User.Name so it is unique
		errPutIndex := rod.PutString(tx, indexUserNameUniqueIndex, userName, userId)
		if errPutIndex != nil {
			return errPutIndex
		}

		fmt.Printf("all done\n")

		return nil
	})

	return &user, err
}
