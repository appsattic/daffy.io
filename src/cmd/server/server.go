package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/gomiddleware/logger"
	"github.com/gomiddleware/mux"
	"github.com/gomiddleware/slash"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/twitter"

	"internal/handlers"
	"internal/middleware"
	"internal/store"
	"internal/types"
)

var sessionName = "session"

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

	// fail if fields haven't been set, or explicitely marked as optional
	valid.SetFieldsRequiredByDefault(true)
}

func main() {
	// setup
	baseUrl := os.Getenv("DAFFY_BASE_URL")
	port := os.Getenv("DAFFY_PORT")
	if port == "" {
		log.Fatal("Specify a port to listen on in the environment variable 'DAFFY_PORT'")
	}
	dbDumpDir := os.Getenv("DAFFY_DB_DUMP_DIR")

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

	// dump this BoltDB to disk every x mins
	if dbDumpDir == "" {
		log.Println("No DB_DUMP_DIR specified - not performing datastore dumps")
	} else {
		ticker := time.NewTicker(1 * time.Hour)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					// do stuff
					log.Println("Dumping the DB now")
					dump(dbDumpDir, boltStore.GetDB())
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

	// Example : https://raw.githubusercontent.com/markbates/goth/master/examples/main.go

	// Twitter
	twitterConsumerKey := os.Getenv("DAFFY_TWITTER_CONSUMER_KEY")
	twitterConsumerSecret := os.Getenv("DAFFY_TWITTER_CONSUMER_SECRET")
	if twitterConsumerKey != "" {
		twitterProvider := twitter.NewAuthenticate(twitterConsumerKey, twitterConsumerSecret, baseUrl+"/auth/twitter/callback")
		goth.UseProviders(twitterProvider)
	}

	// Google (Plus)
	//
	gplusClientId := os.Getenv("DAFFY_GPLUS_CLIENT_ID")
	gplusClientSecret := os.Getenv("DAFFY_GPLUS_CLIENT_SECRET")
	if gplusClientId != "" {
		gplusProvider := gplus.New(gplusClientId, gplusClientSecret, baseUrl+"/auth/gplus/callback")
		goth.UseProviders(gplusProvider)
	}

	// GitHub
	githubClientId := os.Getenv("DAFFY_GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("DAFFY_GITHUB_CLIENT_SECRET")
	if githubClientId != "" {
		githubProvider := github.New(githubClientId, githubClientSecret, baseUrl+"/auth/github/callback")
		goth.UseProviders(githubProvider)
	}

	// Get the providers in use - you could use this to send to your templates so that they know which login links to
	// support, however, you could also just hard-code them in the templates if you're only using one or two.
	providers := goth.GetProviders()
	fmt.Printf("providers=%#v\n", providers)

	// set up some middleware in advance
	checkUser := middleware.CheckUser(sessionStore, sessionName)

	// router
	m := mux.New()

	// do some static routes before doing logging
	m.All("/s", fileServer("static"))
	m.Get("/favicon.ico", serveFile("./static/favicon.ico"))
	m.Get("/robots.txt", serveFile("./static/robots.txt"))

	// some middlewares to always run
	m.Use("/", logger.New())

	// home
	m.Get("/", handlers.HomeHandler(sessionStore, sessionName, providers, tmpl))

	// session
	m.Get("/logout/", slash.Remove)
	m.Get("/logout", handlers.LogoutHandler(sessionStore, sessionName))

	// public user pages
	m.Get("/u", slash.Add)
	m.Get("/u/", redirect("/"))
	m.Get("/u/:username", handlers.ProfileHandler(sessionStore, sessionName, providers, boltStore, tmpl))

	// user routes
	m.Get("/my", slash.Add)
	m.Use("/my", checkUser)
	m.Get("/my/", handlers.MyHandler(sessionStore, sessionName, boltStore, tmpl))

	// tweet from a social account
	m.Get("/my/tweet", middleware.LoadSocials(sessionStore, sessionName, boltStore), handlers.MyTweetHandlerGet(sessionStore, sessionName, boltStore, tmpl))
	m.Post("/my/tweet", middleware.LoadSocials(sessionStore, sessionName, boltStore), handlers.MyTweetHandlerPost(sessionStore, sessionName, boltStore, tmpl))

	// settings
	m.Get("/settings", slash.Add)
	m.Use("/settings", checkUser)
	m.Get("/settings/", handlers.SettingsHandler(sessionStore, sessionName, boltStore, tmpl))
	m.Get("/settings/profile/", slash.Remove)
	m.Post("/settings/profile", handlers.SettingsProfileHandler(sessionStore, sessionName, boltStore))

	// auth
	m.Get("/auth/:provider/", slash.Remove)
	m.Get("/auth/:provider", gothic.BeginAuthHandler)
	m.Get("/auth/:provider/callback", handlers.AuthProviderCallbackHandler(sessionStore, sessionName, boltStore))

	// finally, check all routing was added correctly
	check(m.Err)

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, m)
	check(errServer)
}
