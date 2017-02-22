package handlers

import (
	"html/template"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	"internal/types"
)

func HomeHandler(sessionStore sessions.Store, sessionName string, providers goth.Providers, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("homeHandler"))

		user := getUserFromSession(r, sessionStore, sessionName)

		data := struct {
			Title     string
			User      *types.User
			Providers goth.Providers
		}{
			"Daffy",
			user,
			providers,
		}

		render(w, tmpl, "index.html", data)
	}
}
