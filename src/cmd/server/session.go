package main

import (
	"internal/types"
	"net/http"
)

func getUserFromSession(r *http.Request) *types.User {
	session, _ := sessionStore.Get(r, sessionName)
	user, ok := session.Values["user"].(*types.User)
	if !ok {
		return nil
	}
	return user
}
