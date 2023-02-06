package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// home will write a byte slice and serve as the homepage
func home(w http.ResponseWriter, r *http.Request) {
	// Check to see if the URL from user is anything but a slash
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// Parse the template and check for errors
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Serve the home page template
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// alphaNumeric will check the URL query paramter to ensure only alphanumberic characters are present
func alphaNumeric(name string) (string, bool) {
	return name, regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}

// showCompany handler
func showCompany(w http.ResponseWriter, r *http.Request) {
	// extract the company name parameter from the query string
	name, ok := alphaNumeric(r.URL.Query().Get("name"))
	if !ok || name == "" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific company with name %s...", name)
}

// requestCompany handler
func requestCompany(w http.ResponseWriter, r *http.Request) {
	// Checks if request was made using POST method
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Request a new company here"))
}
