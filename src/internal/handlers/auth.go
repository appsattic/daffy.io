package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gomiddleware/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"

	"internal/store"
)

func AuthProviderCallbackHandler(sessionStore sessions.Store, sessionName string, api store.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("authProviderCallbackHandler"))

		session, _ := sessionStore.Get(r, sessionName)
		currentUser := getUserFromSession(r, sessionStore, sessionName)
		userId := ""
		if currentUser != nil {
			userId = currentUser.Id
		}

		// get this provider name from the router values
		vals := mux.Vals(r)
		provider := vals["provider"]

		authUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("authUser=%#v\n", authUser)

		// check to see if this socialId already exists
		user, err := api.LogIn(userId, provider, authUser.UserID, authUser.NickName, authUser.Name, authUser.Email)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("user=%#v\n", user)

		// we always get a user back from LogIn()
		if user == nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set this info in the session (whether new or updated with a new SocialId)
		session.Values["user"] = &user

		// save all sessions
		sessions.Save(r, w)

		// redirect back to homepage
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
