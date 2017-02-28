package store

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Machiel/slugify"
	valid "github.com/asaskevich/govalidator"
	"github.com/boltdb/bolt"
	"github.com/chilts/rod"
	uuid "github.com/hashicorp/go-uuid"

	"internal/types"
)

var (
	ErrSocialAccountAlreadyExists = errors.New("This social account already exists and is owned by another user.")

	ErrUsernameAlreadyExists = errors.New("Username already exists")
	ErrUsernameUnknown       = errors.New("Unknown username")
)

var userBucket = "user"
var socialBucket = "social"
var indexUserNameUniqueIndex = "i-u-n-u"

type BoltStore struct {
	filename string
	db       *bolt.DB
}

func now() time.Time {
	return time.Now().UTC()
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

func (b *BoltStore) LogIn(userId, provider, id, nickName, title, email string) (*types.User, error) {
	var user types.User
	now := time.Now().UTC()

	isLoggedIn := userId != ""

	// 1. see if this social id exists
	// 2. if it does, read the user and return it
	// 3. if it doesn't, add the Social and User types

	fmt.Printf("boltStore.LogIn(): entry\n")
	fmt.Printf("* userId=%#v\n", userId)
	fmt.Printf("* provider=%#v\n", provider)
	fmt.Printf("* id=%#v\n", id)
	fmt.Printf("* nickName=%#v\n", nickName)
	fmt.Printf("* title=%#v\n", title)
	fmt.Printf("* email=%#v\n", email)

	err := b.db.Update(func(tx *bolt.Tx) error {
		// create a socialId that we use internally (to look the user up)
		socialId := provider + ":" + id

		// fetch this Social entity
		var social types.Social
		errGetSocial := rod.GetJson(tx, socialBucket, socialId, &social)
		if errGetSocial != nil {
			return errGetSocial
		}

		// check to see if the socialId exists
		if social.Id != "" {
			fmt.Printf("Got social = %#v\n", social)

			// Now that we have the user (from an existing Social), check to see if this matches the current user if
			// one is currently logged in.
			if isLoggedIn {
				if social.UserId != userId {
					fmt.Printf("Social account already exists and is owned by another user already")
					return ErrSocialAccountAlreadyExists
				}
			}

			// get this user - should ALWAYS work if the above Social exists
			errGetUser := rod.GetJson(tx, userBucket, social.UserId, &user)
			fmt.Printf("Got user = %#v\n", user)
			if errGetUser != nil {
				return errGetUser
			}
			return nil
		}

		// check to see if we already have a user logged in
		if !isLoggedIn {
			// no user, so create a unique UserId for this user (this never changes)
			userId, _ = uuid.GenerateUUID()
		}

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

		// if we have a user logged in already, get the user and add this socialId
		var userName string
		if isLoggedIn {
			errGetJson := rod.GetJson(tx, userBucket, userId, &user)
			if errGetJson != nil {
				return errGetJson
			}

			user.SocialIds = append(user.SocialIds, socialId)
			user.Updated = now
			userName = user.Name
		} else {
			// create a unique userName for this user - they can change it if they like
			userName = slugify.Slugify(nickName + "-" + id)

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
		}

		fmt.Printf("* created userId=%#v\n", userId)
		errPutUser := rod.PutJson(tx, userBucket, userId, user)
		if errPutUser != nil {
			return errPutUser
		}

		// index the user.Name
		if isLoggedIn {
			// no need to add this userName to the index since it is already there
		} else {
			// index the User.Name so it is unique
			errPutIndex := rod.PutString(tx, indexUserNameUniqueIndex, userName, userId)
			if errPutIndex != nil {
				return errPutIndex
			}
		}

		fmt.Printf("all done\n")

		return nil
	})

	return &user, err
}

func (b *BoltStore) SelSocials(socialIds []string) ([]types.Social, error) {
	socials := make([]types.Social, 0)

	err := b.db.View(func(tx *bolt.Tx) error {
		for _, socialId := range socialIds {
			social := types.Social{}

			errGetJson := rod.GetJson(tx, socialBucket, socialId, &social)
			if errGetJson != nil {
				return errGetJson
			}

			socials = append(socials, social)
		}

		return nil
	})

	return socials, err
}

func (b *BoltStore) GetUserPublic(username string) (*types.User, error) {
	var user *types.User

	err := b.db.View(func(tx *bolt.Tx) error {
		// firstly, read the index
		userId, errGetIndex := rod.GetString(tx, indexUserNameUniqueIndex, username)
		if errGetIndex != nil {
			return errGetIndex
		}
		if userId == "" {
			return nil
		}

		// now get this user
		var newUser types.User
		errGetJson := rod.GetJson(tx, userBucket, userId, &newUser)
		if errGetJson != nil {
			return errGetJson
		}
		user = &newUser
		return nil
	})

	return user, err
}

func (b *BoltStore) UpdateUser(currentUser types.User, updateUser types.UpdateUser) (types.User, error) {
	var user types.User
	now := now()

	// first thing to do is validate the incoming info
	isValid, errs := valid.ValidateStruct(updateUser)
	if errs != nil {
		return user, errs
	}

	log.Printf("isValid = %#v\n", isValid)
	log.Printf("errs    = %#v\n", errs)

	if !isValid {
		// ToDo: ... !!!

		// re-render the page with the errors

		// ToDo: ... !!!
	}

	// Steps:
	// 1. Get the full user out of the store
	// 2. Update with the new fields
	// 3. Re-save back out

	err := b.db.Update(func(tx *bolt.Tx) error {
		// get this user - should ALWAYS work if the above Social exists
		errGetUser := rod.GetJson(tx, userBucket, currentUser.Id, &user)
		fmt.Printf("Got user = %#v\n", user)
		if errGetUser != nil {
			return errGetUser
		}

		// check to see if the username has changed, and if so, remove the old index entry and add a new one
		if updateUser.Name != user.Name {
			fmt.Printf("User is changing their username from %s, to %s\n", user.Name, updateUser.Name)

			// check that this username doesn't already exist
			id, errGetIndex := rod.GetString(tx, indexUserNameUniqueIndex, updateUser.Name)
			if errGetIndex != nil {
				return errGetIndex
			}
			if id != "" {
				return ErrUsernameAlreadyExists
			}

			// remove the old index value
			errDel := rod.Del(tx, indexUserNameUniqueIndex, user.Name)
			if errDel != nil {
				return errDel
			}

			// put the new index value
			errPutIndex := rod.PutString(tx, indexUserNameUniqueIndex, updateUser.Name, currentUser.Id)
			if errPutIndex != nil {
				return errPutIndex
			}
		} else {
			fmt.Printf("User is NOT changing their username.\n")
		}

		// update
		user.Name = updateUser.Name
		user.Title = updateUser.Title
		user.Email = updateUser.Email
		user.Updated = now

		// re-save user
		errPutUser := rod.PutJson(tx, userBucket, user.Id, user)
		fmt.Printf("errPutUser = %#v\n", errPutUser)
		if errPutUser != nil {
			return errPutUser
		}

		return nil
	})

	return user, err
}

func (b *BoltStore) GetDB() *bolt.DB {
	return b.db
}
