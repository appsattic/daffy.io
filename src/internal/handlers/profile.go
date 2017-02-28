package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gomiddleware/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	"internal/store"
	"internal/types"
)

func ProfileHandler(sessionStore sessions.Store, sessionName string, providers goth.Providers, boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("UserHandler"))

		user := getUserFromSession(r, sessionStore, sessionName)

		vals := mux.Vals(r)
		fmt.Printf("username=%s\n", vals["username"])

		// get this user from the store
		profile, err := boltStore.GetUserPublic(vals["username"])
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if profile == nil {
			http.NotFound(w, r)
			return
		}

		data := struct {
			Title     string
			User      *types.User
			Providers goth.Providers
			Profile   *types.User
		}{
			"User Profile - daffy.io",
			user,
			providers,
			profile,
		}
		render(w, tmpl, "u-:username.html", data)
	}
}
