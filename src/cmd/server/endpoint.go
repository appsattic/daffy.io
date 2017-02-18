package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomiddleware/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"

	"internal/store"
	"internal/types"
)

var sessionName = "session"

// create a decoder than can be used for all forms
var decoder = schema.NewDecoder()

// To create newer keys, setup two V3 environment variables and drop the V1 ones (or keep them for a while). Eventually
// you can drop them. Keep incrementing each time you add new ones. See : https://godoc.org/github.com/gorilla/sessions
var sessionStore = sessions.NewCookieStore(
	// New Keys
	[]byte(os.Getenv("DAFFY_SESSION_AUTH_KEY_V2")),
	[]byte(os.Getenv("DAFFY_SESSION_ENC_KEY_V2")),
	// Old Keys
	[]byte(os.Getenv("DAFFY_SESSION_AUTH_KEY_V1")),
	[]byte(os.Getenv("DAFFY_SESSION_ENC_KEY_V1")),
)

func init() {
	// tell gothic where our session store is
	gothic.Store = sessionStore
}

func homeHandler(tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)
		user := getUserFromSession(session)

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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionName)

	// scrub user
	delete(session.Values, "user")
	session.Save(r, w)

	// redirect to somewhere else
	http.Redirect(w, r, "/", http.StatusFound)
}

func settingsProfileHandler(boltStore *store.BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// firstly, check the user is logged in
		session, _ := sessionStore.Get(r, sessionName)
		user := getUserFromSession(session)
		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

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
		session.Values["user"] = &newUser
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func myHandler(boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check the user is logged in
		session, _ := sessionStore.Get(r, sessionName)
		user := getUserFromSession(session)
		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

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

func settingsHandler(boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)
		user := getUserFromSession(session)
		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

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

func authProviderCallbackHandler(boltStore *store.BoltStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)
		currentUser := getUserFromSession(session)
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
		user, err := boltStore.LogIn(userId, provider, authUser.UserID, authUser.NickName, authUser.Name, authUser.Email)
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
