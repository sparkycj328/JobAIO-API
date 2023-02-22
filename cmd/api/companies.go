package main

import (
	"fmt"
	"github.com/sparkycj328/JobAIO-API/internal/data"
	"net/http"
)

// createCompanyHandler will insert job postings into the database based on the company name
func (app *application) createCompanyHandler(w http.ResponseWriter, r *http.Request) {

	// create a local copy of the company struct which will store the request body
	var input struct {
		Name      string         `json:"company"`   // company name
		Countries []data.Country `json:"countries"` // job postings
	}

	// decode the request body
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v", input)
}

// showCompanyHandler will display the job information for the specified company
func (app *application) showCompanyHandler(w http.ResponseWriter, r *http.Request) {
	name, err := app.readNameParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// otherwise interpolate the name parameter
	fmt.Fprintf(w, "Show the details of %s", name)
}
