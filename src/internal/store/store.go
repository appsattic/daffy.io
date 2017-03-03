package store

import (
	"internal/types"

	"github.com/markbates/goth"
)

type Api interface {
	Open() error
	Close() error

	// socialId is just "twitter-123456", "facebook-777", or "github-13579"
	LogInGoth(userId, provider string, authUser goth.User) (*types.User, error)
	LogIn(userId, provider, socialId, socialUserName, title, email, accessToken, accessTokenSecret string) (*types.User, error)
	SelSocials(socialIds []string) ([]types.Social, error) // ToDo: check if this should be in the API

	// The following API are public and don't require a `currentUser`.
	GetUserPublic(username string) (*types.User, error)

	// The following API calls require a `currentUser` so we know the user is authenticated.
	UpdateUser(currentUser types.User, data types.UpdateUser) (types.User, error)
}
