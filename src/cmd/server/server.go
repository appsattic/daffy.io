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
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
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

	// create a decoder than can be used for all forms
	var decoder = schema.NewDecoder()

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

	m.All("/s", fileServer("static"))
	m.Get("/favicon.ico", serveFile("./static/favicon.ico"))
	m.Get("/robots.txt", serveFile("./static/robots.txt"))

	// some middlewares to always run
	m.Use("/", logger.New())

	m.Get("/auth/:provider/callback", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// begin auth
	m.Get("/auth/:provider", gothic.BeginAuthHandler)

	// logout
	m.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, sessionName)

		// scrub user
		delete(session.Values, "user")
		session.Save(r, w)

		// redirect to somewhere else
		http.Redirect(w, r, "/", http.StatusFound)
	})

	m.Post("/settings/profile", func(w http.ResponseWriter, r *http.Request) {
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
	})

	m.Get("/my", slash.Add)
	m.Get("/my/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	m.Get("/settings/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
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

	check(m.Err)

	// server
	log.Printf("Starting server, listening on port %s\n", port)
	errServer := http.ListenAndServe(":"+port, m)
	check(errServer)
}
