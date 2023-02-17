package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

// alphaNumeric will check the URL query paramter to ensure only alphanumberic characters are present
func alphaNumeric(name string) (string, bool) {
	return name, regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}

// showCompanyHandler will display the job information for the specified company
func (app *application) showCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve a slice containing any interpolated parameter names and values
	params := httprouter.ParamsFromContext(r.Context())

	// retrieve the name parameter
	name, ok := alphaNumeric(params.ByName("name"))
	if !ok || name == "" {
		http.NotFound(w, r)
		return
	}
	// otherwise interpolate the name parameter
	fmt.Fprintf(w, "Show the details of %s", name)
}
