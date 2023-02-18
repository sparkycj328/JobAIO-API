package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// routes returns a httprouter.Router and sets the appropriate endpoints with each method
func (app *application) routes() *httprouter.Router {

	// Initialize a new httprouter interface
	router := httprouter.New()

	// register the appropriate methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/company/:name", app.showCompanyHandler)

	return router
}
