package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/moniquelin/monday-hr/internal/data"
)

// authenticate checks whether a user is logged in
func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check Authorization header for JWT token
		// 1st check: check if empty
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		tokenString = tokenString[len("Bearer "):]

		// 2nd check: verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.Config.Jwt.Secret), nil
		})

		if err != nil || !token.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get user ID from claims
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			app.invalidAuthenticationTokenResponse(w, r)
			fmt.Print("test 3")
			return
		}

		// Retrieve the details of the user associated with the authentication token
		user, err := app.Models.Users.Get(int64(userIDFloat))
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
				fmt.Print("test 4")

			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Call contextSetUser() to add the user information to the request context.
		r = app.contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// requireAdmin checks whether the given user has the "admin" role
func (app *Application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.Role == "employee" {
			app.errorResponse(w, r, http.StatusUnauthorized, "you must be an admin to access this resource")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireEmployee checks whether the given user has the "employee" role
func (app *Application) requireEmployee(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.Role != "employee" {
			app.errorResponse(w, r, http.StatusUnauthorized, "you must be an employee to access this resource")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireAllowSeed checks if the MONDAY_HR_ALLOW_SEED env var is "true"
func (app *Application) requireAllowSeed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowSeed := os.Getenv("MONDAY_HR_ALLOW_SEED")
		switch allowSeed {
		case "":
			app.errorResponse(w, r, http.StatusUnauthorized, "MONDAY_HR_ALLOW_SEED is not set!")
		case "true":
			next.ServeHTTP(w, r)
		default:
			app.errorResponse(w, r, 403, "You have no access to seed data!")
			return
		}
	})
}
