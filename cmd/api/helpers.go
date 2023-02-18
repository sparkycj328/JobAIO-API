package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

// alphaNumeric will check the URL query paramter to ensure only alphanumberic characters are present
func alphaNumeric(name string) (string, bool) {
	return name, regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}

// readNameParam will grab the parameter from the URL using the request context
func (app *application) readNameParam(r *http.Request) (string, error) {
	// retrieve a slice containing any interpolated parameter names and values
	params := httprouter.ParamsFromContext(r.Context())

	// retrieve the name parameter
	name, ok := alphaNumeric(params.ByName("name"))
	if !ok || name == "" {
		return "", errors.New("invalid name parameter")
	}

	return name, nil
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
