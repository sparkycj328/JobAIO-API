package main

import (
	"fmt"
	"net/http"
)

// healthcheckHandler writes plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	// create a JSON-object to return to the response writer
	js := `{"status":"available", "environment":%q, "version":%q}`
	js = fmt.Sprintf(js, app.config.env, version)

	// sets a response header so client knows response contains json
	w.Header().Set("Content-Type", "application/json")

	// write the JSON as the HTTP response body
	w.Write([]byte(js))

}
