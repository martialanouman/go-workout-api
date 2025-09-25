package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/martialanouman/femProject/internal/store"
	"github.com/martialanouman/femProject/internal/tokens"
	"github.com/martialanouman/femProject/internal/utils"
)

type UserMiddleware struct {
	Store store.UserStore
}

type UserContextKey string

const (
	UserContextKeyName = UserContextKey("user")
	BearerTokenName    = "Bearer"
)

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKeyName, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKeyName).(*store.User)
	if !ok {
		panic("could not get user from context")
	}

	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := w.Header().Get("Authorization")

		if authHeader == "" {
			r := SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != BearerTokenName {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or missing authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.Store.GetUserByToken(tokens.ScopeAuth, token)
		if err != nil || user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired auth token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
