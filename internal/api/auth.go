package api

import (
	"errors"
	"net/http"

	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/validator"
)

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
	tokenString, err := app.createToken(user.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.writeJSON(w, 200, envelope{"user": user, "authentication_token": tokenString}, nil)
}
