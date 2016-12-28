package main

import (
	"github.com/gorilla/sessions"

	"internal/types"
)

func getUserFromSession(session *sessions.Session) *types.User {
	user, ok := session.Values["user"].(*types.User)
	if !ok {
		return nil
	}
	return user
}
