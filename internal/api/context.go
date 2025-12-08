package api

import (
	"context"
	"net/http"

	"github.com/moniquelin/monday-hr/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context
func (app *Application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextGetUser() retrieves the User struct from the request context
func (app *Application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
