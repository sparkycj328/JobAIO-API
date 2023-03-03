package main

import (
	"errors"
	"fmt"
	"github.com/sparkycj328/JobAIO-API/internal/data"
	"github.com/sparkycj328/JobAIO-API/internal/validator"
	"net/http"
	"time"
)

// createCompanyHandler will insert job postings into the database based on the company name
func (app *application) createCompanyHandler(w http.ResponseWriter, r *http.Request) {

	// create a local copy of the company struct which will store the request body
	var input struct {
		ID        int64      `json:"-"`                 // Unique integer id for the company
		Name      string     `json:"company"`           // company name
		Country   string     `json:"country"`           // Country name
		Total     int        `json:"total"`             // total amount of job available
		URL       string     `json:"url"`               // URL location where resource is located
		CreatedAt *time.Time `json:"created,omitempty"` // created timestamp for the data
	}

	// decode the request body
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	company := &data.Company{
		Name:    input.Name,
		Country: input.Country,
		Total:   input.Total,
		URL:     input.URL,
	}

	// Initialize a new validator
	v := validator.New()
	if data.ValidateCompany(v, company); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the data into the jobs table
	if err = app.models.Vendors.Insert(company); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Custom header declaration in order to pass the location of the records
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/companies/%s", company.Name))

	// Write a JSON response with a 201 created status code, the vendor data in the response
	// body and the Location folder
	if err := app.writeJSON(w, http.StatusCreated, envelope{"company": company}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// getRecordHandler will execute our single row query based on the id
// parameter which is grabbed from the context from the request
func (app *application) getRecordHandler(w http.ResponseWriter, r *http.Request) {

}

// showCompanyHandler will display the job information for the specified company
func (app *application) showCompanyHandler(w http.ResponseWriter, r *http.Request) {
	name, err := app.readNameParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	jobs, err := app.models.Vendors.GetRows(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"jobs": jobs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
