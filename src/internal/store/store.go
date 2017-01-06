package store

import (
	"internal/types"
)

type Api interface {
	Open() error
	Close() error

	// socialId is just "twitter-123456", "facebook-777", or "github-13579"
	LogIn(userId, provider, socialId, socialUserName, title, email string) (*types.User, error)
}
