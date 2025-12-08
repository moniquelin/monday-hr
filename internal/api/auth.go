package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/validator"
)

// loginHandler manages user login
func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email             string `json:"email"`
		PlaintextPassword string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate Email & Password
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.PlaintextPassword)
	if len(v.Errors) != 0 {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get user by email from the database
	var user *data.User

	user, err = app.Models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		// If no user is found
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, 401, "user with this email does not exist")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Match the input password with the user password
	var passwordMatch bool
	passwordMatch, err = user.Password.Matches(input.PlaintextPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !passwordMatch {
		app.errorResponse(w, r, 401, "wrong password!")
		return
	}

	// Create token for user with correct email and password
	tokenString, err := app.createToken(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.writeJSON(w, 201, envelope{"user": user, "authentication_token": tokenString}, nil)
}

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
			return app.Config.Jwt.Secret, nil
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

		fmt.Println(claims)

		// Get user ID from claims
		userID, ok := claims["sub"].(int64)
		if !ok {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Retrieve the details of the user associated with the authentication token
		user, err := app.Models.Users.Get(userID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
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
