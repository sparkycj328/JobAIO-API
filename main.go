package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// home will write a byte slice
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
func main() {

	// create a new servemux and register the home function to act as the handler
	// for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/company", showCompany)
	mux.HandleFunc("/company/request", requestCompany)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
