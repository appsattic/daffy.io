package handlers

import (
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gorilla/sessions"

	"internal/sess"
	"internal/types"
)

func getUserFromSession(r *http.Request, sessionStore sessions.Store, sessionName string) *types.User {
	return sess.GetUserFromSession(r, sessionStore, sessionName)
}

func LogoutHandler(sessionStore sessions.Store, sessionName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("logoutHandler"))

		session, _ := sessionStore.Get(r, sessionName)

		// scrub user
		delete(session.Values, "user")
		session.Save(r, w)

		// redirect to somewhere else
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
