package middleware

import (
	"log"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gorilla/sessions"

	"internal/types"
)

func CheckUser(sessionStore sessions.Store, sessionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer logfn.Exit(logfn.Enter("main.checkUser"))

			// get the session
			session, _ := sessionStore.Get(r, sessionName)

			// now check that the "user" key exists
			value, ok := session.Values["user"]
			if !ok {
				log.Println("main.checkUser(): no user key in session")
				// no "user" key
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			// assert the value is user
			_, ok = value.(*types.User)
			if !ok {
				log.Println("main.checkUser(): user key in session can't assert to a valid user")
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			// serve the next middleware
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
