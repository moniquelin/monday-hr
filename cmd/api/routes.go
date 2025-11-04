package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance
	router := httprouter.New()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints
	router.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)

	// Return the httprouter instance
	return router
}
