package sess

import (
	"net/http"

	"github.com/gorilla/sessions"

	"internal/types"
)

func GetUserFromSession(r *http.Request, sessionStore sessions.Store, sessionName string) *types.User {
	session, _ := sessionStore.Get(r, sessionName)
	user, ok := session.Values["user"].(*types.User)
	if !ok {
		return nil
	}
	return user
}
