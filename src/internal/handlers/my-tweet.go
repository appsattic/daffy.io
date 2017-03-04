package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chilts/logfn"
	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"

	"internal/middleware"
	"internal/store"
	"internal/types"
)

func MyTweetHandlerGet(sessionStore sessions.Store, sessionName string, boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("handlers.MyTweetHandlerGet"))

		user := getUserFromSession(r, sessionStore, sessionName)
		socials := middleware.GetSocials(r)

		data := struct {
			Title   string
			User    *types.User
			Socials []types.Social
		}{
			"Tweets - daffy.io",
			user,
			socials,
		}
		render(w, tmpl, "my-tweet.html", data)
	}
}

func MyTweetHandlerPost(sessionStore sessions.Store, sessionName string, boltStore *store.BoltStore, tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer logfn.Exit(logfn.Enter("handlers.MyTweetHandlerPost"))

		socialId := r.FormValue("SocialId")
		tweet := r.FormValue("Tweet")
		fmt.Printf("socialid=%#v\n", socialId)
		fmt.Printf("tweet=%#v\n", tweet)

		// ToDo: check that there is something in the tweet, and figure out the 140 char rules.
		// ToDo: check that socialId is contained within the socials list (or actually, just check for ourselves in the store)
		socials := middleware.GetSocials(r)

		// let's post this tweet on behalf of the user
		consumerKey := os.Getenv("DAFFY_TWITTER_CONSUMER_KEY")
		consumerSecret := os.Getenv("DAFFY_TWITTER_CONSUMER_SECRET")
		consumer := oauth.NewConsumer(consumerKey, consumerSecret, oauth.ServiceProvider{})
		// consumer.Debug(true)

		// create the accessToken
		accessToken := &oauth.AccessToken{
			Token:  socials[0].AccessToken,
			Secret: socials[0].AccessTokenSecret,
		}
		twitterEndPoint := "https://api.twitter.com/1.1/statuses/update.json"
		params := make(map[string]string)
		params["status"] = tweet
		response, err := consumer.Post(twitterEndPoint, params, accessToken)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		// get the response and perhaps tell the user about it
		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(string(respBody))

		http.Redirect(w, r, "/my/", http.StatusFound)
	}
}
