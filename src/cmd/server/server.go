package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomiddleware/logger"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
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

func render(w http.ResponseWriter, tmpl *template.Template, tmplName string, data interface{}) {
	buf := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
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
	if port == "" {
		log.Fatal("Specify a port to listen on in the environment variable 'DAFFY_PORT'")
	}

	// load up all templates
	tmpl, err := template.New("").ParseGlob("./templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// create/open/connect to a store
	boltStore := store.NewBoltStore("daffy.db")
	errOpen := boltStore.Open()
	check(errOpen)
	defer boltStore.Close()

	// Example : https://raw.githubusercontent.com/markbates/goth/master/examples/main.go

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
	twitterProvider := twitter.NewAuthenticate(twitterConsumerKey, twitterConsumerSecret, baseUrl+"/auth/twitter/callback")

	// GitHub
	//
	// Follow the instructions here or here:
	//
	// * https://github.com/settings/developers
	// * https://github.com/organizations/<your-organization>/settings/applications
	//
	// Requires:
	//
	// * DAFFY_GITHUB_CLIENT_ID
	// * DAFFY_GITHUB_CLIENT_SECRET
	githubClientId := os.Getenv("DAFFY_GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("DAFFY_GITHUB_CLIENT_SECRET")
	githubProvider := github.New(githubClientId, githubClientSecret, baseUrl+"/auth/github/callback")

	// goth
	goth.UseProviders(twitterProvider, githubProvider)

	// router
	p := pat.New()

	p.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)
		currentUser := getUserFromSession(session)
		userId := ""
		if currentUser != nil {
			userId = currentUser.Id
		}

		// get this provider name from the URL
		provider := r.URL.Query().Get(":provider")

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

		data := struct {
			Title string
			User  *types.User
		}{
			"Daffy",
			user,
		}

		render(w, tmpl, "index.html", data)
	})

	// create the logger middleware
	lgr := logger.New()

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, lgr(p))
	check(errServer)
}
