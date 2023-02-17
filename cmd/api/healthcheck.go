package main

import (
	"encoding/json"
	"net/http"
)

// healthcheckHandler writes plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	// marshal the json object into bytes
	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request",
			http.StatusInternalServerError)
		return
	}

	// append a newline to make the JSON more viewer-friendly
	js = append(js, '\n')

	// sets a response header so client knows response contains json
	w.Header().Set("Content-Type", "application/json")

	// write the JSON as the HTTP response body
	w.Write(js)
}
