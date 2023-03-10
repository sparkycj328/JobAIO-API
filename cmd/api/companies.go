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

// showRecordHandler will execute our single row query based on the id
// parameter which is grabbed from the context from the request
func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	// grab the id parameter from the request url
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	record, err := app.models.Vendors.GetRecord(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// if no error was returned by our Get record query, write the record to JSON
	if err := app.writeJSON(w, http.StatusOK, envelope{"record": record}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listCompanyHandler will display the job information for the specified company
// the api endpoint can receive several query parameters and will filter our table
// based on these query parameters
func (app *application) listCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// create a local copy of the company struct which will store the request body
	var input struct {
		Name    string
		Country string
		Total   int
		data.Filters
	}
	// initialize a new validator struct
	v := validator.New()

	// Retrieve the url query parameter map from url.Values
	qs := r.URL.Query()

	// Retrieve the url query values and store them in our struct
	input.Name = app.readString(qs, "vendor", "")
	input.Country = app.readString(qs, "country", "")
	input.Total = app.readInt(qs, "total", 0, v)

	// get the page and page_size query string values as integers and set
	// their default values
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 50, v)

	// Extract the sort query parameter to sort the database query by
	// default sort if no sort parameter is provided will sort the query
	// by the vendor name in ascending order
	input.Filters.SortSafeList = []string{"id", "vendor", "amount", "created_at", "country",
		"-id", "-vendor", "-amount", "-created_at", "-country"}
	input.Filters.Sort = app.readString(qs, "sort", "vendor")

	// Check the Validator error map for any errors added by our app.readInt method
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	// Call the GetAll function in order to grab all rows
	jobs, err := app.models.Vendors.GetAllRows(input.Name, input.Country, input.Total, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"jobs": jobs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
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

	if err := app.writeJSON(w, http.StatusOK, envelope{"jobs": &jobs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateCompanyHandler will update a record based on the ID parameter
func (app *application) updateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve the id parameter
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	// fetch the individual record to be updated
	record, err := app.models.Vendors.GetRecord(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// create a local copy of the company struct which will store the request body
	var input struct {
		ID        int64      `json:"-"`                 // Unique integer id for the company
		Name      string     `json:"company"`           // company name
		Country   string     `json:"country"`           // Country name
		Total     int        `json:"total"`             // total amount of job available
		URL       string     `json:"url"`               // URL location where resource is located
		CreatedAt *time.Time `json:"created,omitempty"` // created timestamp for the data
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
	}

	// read the json data from the local input struct into the returned record struct
	record.Name = input.Name
	record.Country = input.Country
	record.Total = input.Total
	record.URL = input.URL

	// validate that the json data is valid before updating the record in our table
	v := validator.New()
	if data.ValidateCompany(v, record); v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	// write the new company struct to our database
	if err := app.models.Vendors.Update(record); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// once updated, return the JSON struct to the client making the request
	if err := app.writeJSON(w, http.StatusOK, envelope{"updated": record}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteCompany will delete a single record
func (app *application) deleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// grab the id parameter from the url
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	// pass the id parameter to the delete function
	err = app.models.Vendors.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// write a JSON response upon successful deletion of the record
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "record successfully updated"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
