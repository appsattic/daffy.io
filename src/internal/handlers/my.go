package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	"internal/store"
	"internal/types"
)

// create a decoder than can be used for all forms
var decoder = schema.NewDecoder()

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

func SettingsProfileHandler(sessionStore sessions.Store, sessionName string, boltStore *store.BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("settingsProfileHandler"))

		user := getUserFromSession(r, sessionStore, sessionName)

		// parse the incoming form
		err := r.ParseForm()
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// decode the form into a types.UpdateUser
		updateUser := types.UpdateUser{}
		err = decoder.Decode(&updateUser, r.PostForm)
		// check if this errors is from `govalidator` rather than any other general kind of error
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// tell the API what to update
		fmt.Printf("updateUser=%#v\n", updateUser)

		// update this user
		newUser, err := boltStore.UpdateUser(*user, updateUser)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("* NEW * user = %v\n", newUser)

		// save this new user
		session, _ := sessionStore.Get(r, sessionName)
		session.Values["user"] = &newUser
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func MyHandler(sessionStore sessions.Store, sessionName string, boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("myHandler"))

		user := getUserFromSession(r, sessionStore, sessionName)

		// get all the social entities
		socials, err := boltStore.SelSocials(user.SocialIds)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			User    *types.User
			Socials []types.Social
		}{
			"My Daffy - daffy.io",
			user,
			socials,
		}
		render(w, tmpl, "my-index.html", data)
	}
}

func SettingsHandler(sessionStore sessions.Store, sessionName string, boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("settingsHandler"))

		user := getUserFromSession(r, sessionStore, sessionName)

		// get all the social entities
		socials, err := boltStore.SelSocials(user.SocialIds)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			User    *types.User
			Socials []types.Social
		}{
			"Settings - daffy.io",
			user,
			socials,
		}
		render(w, tmpl, "settings-index.html", data)
	}
}
