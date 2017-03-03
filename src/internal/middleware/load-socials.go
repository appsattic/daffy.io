package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/chilts/logfn"
	"github.com/gorilla/sessions"

	"internal/sess"
	"internal/store"
	"internal/types"
)

const socialsKey key = 42

func LoadSocials(sessionStore sessions.Store, sessionName string, boltStore *store.BoltStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer logfn.Exit(logfn.Enter("handlers.CheckUser"))

			user := sess.GetUserFromSession(r, sessionStore, sessionName)

			// get all the social entities
			socials, err := boltStore.SelSocials(user.SocialIds)
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// store this in the context
			ctx := context.WithValue(r.Context(), socialsKey, socials)

			// serve the next middleware
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

// GetSocials can be used to obtain the Socials related to this user.
func GetSocials(r *http.Request) []types.Social {
	return r.Context().Value(socialsKey).([]types.Social)
}
