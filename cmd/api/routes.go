package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// routes returns a httprouter.Router and sets the appropriate endpoints with each method
func (app *application) routes() http.Handler {

	// Initialize a new httprouter interface
	router := httprouter.New()

	// change httprouters error handling to app's error handling for errors 404 and 405
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// register the appropriate methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/companies", app.listCompanyHandler)
	router.HandlerFunc(http.MethodPost, "/v1/companies", app.createCompanyHandler)
	router.HandlerFunc(http.MethodGet, "/v1/companies/:name", app.showCompanyHandler)
	router.HandlerFunc(http.MethodPut, "/v1/companies/:id", app.updateCompanyHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/companies/:id", app.deleteCompanyHandler)
	router.HandlerFunc(http.MethodGet, "/v1/record/:id", app.showRecordHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
