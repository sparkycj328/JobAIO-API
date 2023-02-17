package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

// alphaNumeric will check the URL query paramter to ensure only alphanumberic characters are present
func alphaNumeric(name string) (string, bool) {
	return name, regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}

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
