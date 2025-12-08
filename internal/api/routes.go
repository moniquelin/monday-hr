package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()

	// Public routes
	router.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", app.loginHandler)

	// Protected routes (Employee Only)
	router.Handler(http.MethodPost, "/v1/attendance/checkin",
		app.authenticate(app.requireEmployee(http.HandlerFunc(app.checkInHandler))))

	return router
}
