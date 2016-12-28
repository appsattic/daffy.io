package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"

	"internal/store"
	"internal/types"
)

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

var sessionName = "session"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// tell gothic where our session store is
	gothic.Store = sessionStore

	// Register the user with `gob` so we can serialise it.
	gob.Register(&types.User{})
}

func main() {
	// setup
	baseUrl := os.Getenv("DAFFY_BASE_URL")
	port := os.Getenv("DAFFY_PORT")

	// create/open/connect to a store
	boltStore := store.NewBoltStore("daffy.db")
	errOpen := boltStore.Open()
	check(errOpen)
	defer boltStore.Close()

	// Twitter
	//
	// Create a new app : https://apps.twitter.com/app/new
	//
	// Requires:
	//
	// * DAFFY_TWITTER_CONSUMER_KEY
	// * DAFFY_TWITTER_CONSUMER_SECRET
	twitterConsumerKey := os.Getenv("DAFFY_TWITTER_CONSUMER_KEY")
	twitterConsumerSecret := os.Getenv("DAFFY_TWITTER_CONSUMER_SECRET")
	twitter := twitter.NewAuthenticate(twitterConsumerKey, twitterConsumerSecret, baseUrl+"/auth/twitter/callback")

	// goth
	goth.UseProviders(twitter)

	// router
	p := pat.New()

	p.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)

		// get this provider name from the URL
		provider := r.URL.Query().Get(":provider")

		authUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("authUser=%#v\n", authUser)

		// ToDo: check in the datastore for this Social ID - insert if it doesn't exist.

		// for now, just create a User that we can store in the session
		user := types.User{
			Name:  provider + "-" + authUser.UserID + "-" + authUser.NickName,
			Title: authUser.Name,
			Email: authUser.Email,
		}

		// set this info in the session
		session.Values["user"] = &user

		// save all sessions
		sessions.Save(r, w)

		// redirect back to homepage
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// begin auth
	p.Get("/auth/{provider}", gothic.BeginAuthHandler)

	// logout
	p.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)

		// scrub user
		delete(session.Values, "user")
		session.Save(r, w)

		// redirect to somewhere else
		http.Redirect(w, r, "/", http.StatusFound)
	})

	p.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		session, _ := sessionStore.Get(r, sessionName)
		user := getUserFromSession(session)

		if user != nil {
			// a session
			fmt.Fprintf(w, "<p>You are logged in as %s. <a href='/logout'>Log Out.</a></p>", user.Name)
		} else {
			// no session
			fmt.Fprintf(w, "<p>You are not logged in. <a href='/auth/twitter'>Log in with Twitter.</a></p>")
		}
	})

	// server
	errServer := http.ListenAndServe(":"+port, p)
	check(errServer)
}
