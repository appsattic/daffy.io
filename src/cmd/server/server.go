package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	valid "github.com/asaskevich/govalidator"
	"github.com/gomiddleware/logger"
	"github.com/gomiddleware/mux"
	"github.com/gomiddleware/slash"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/twitter"

	"internal/store"
	"internal/types"
)

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

	// router
	m := mux.New()

	// do some static routes before doing logging
	m.All("/s", fileServer("static"))
	m.Get("/favicon.ico", serveFile("./static/favicon.ico"))
	m.Get("/robots.txt", serveFile("./static/robots.txt"))

	// some middlewares to always run
	m.Use("/", logger.New())

	// home
	m.Get("/", homeHandler(tmpl))

	// session
	m.Get("/logout/", slash.Remove)
	m.Get("/logout", logoutHandler)

	// user routes
	m.Get("/my", slash.Add)
	m.Get("/my/", myHandler(boltStore, tmpl))
	m.Get("/settings", slash.Add)
	m.Get("/settings/", settingsHandler(boltStore, tmpl))
	m.Get("/settings/profile/", slash.Remove)
	m.Post("/settings/profile", settingsProfileHandler(boltStore))

	// auth
	m.Get("/auth/:provider/", slash.Remove)
	m.Get("/auth/:provider", gothic.BeginAuthHandler)
	m.Get("/auth/:provider/callback", authProviderCallbackHandler(boltStore))

	// finally, check all routing was added correctly
	check(m.Err)

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, m)
	check(errServer)
}
