package main

import (
	"encoding/json"
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

// writeJSON is a helper function which will iterate through the data passed and convert it into
// a JSON object to return.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// append a newline to make the JSON more viewer-friendly
	js = append(js, '\n')

	// iterate through map and write the Header for each key-value pair
	for key, value := range headers {
		w.Header()[key] = value
	}

	// sets a response header so client knows response contains json
	w.Header().Set("Content-Type", "application/json")

	// write the status code for the response
	w.WriteHeader(status)
	// write the JSON as the HTTP response body
	w.Write(js)

	return nil
}
