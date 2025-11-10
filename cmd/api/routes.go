package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return router
}
