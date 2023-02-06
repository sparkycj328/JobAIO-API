package main

import (
	"fmt"
	"net/http"
	"regexp"
)

// home will write a byte slice and serve as the homepage
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from JobAIO"))
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
