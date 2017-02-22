package handlers

import (
	"html/template"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gorilla/sessions"

	"internal/types"
)

func HomeHandler(sessionStore sessions.Store, sessionName string, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("homeHandler"))

		user := getUserFromSession(r, sessionStore, sessionName)

		data := struct {
			Title string
			User  *types.User
		}{
			"Daffy",
			user,
		}

		render(w, tmpl, "index.html", data)
	}
}
