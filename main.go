package main

import (
	"log"
	"net/http"
)

// home will write a byte slice
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

// showCompany handler
func showCompany(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display jobs for a specific company"))
}

// requestCompany handler
func requestCompany(w http.ResponseWriter, r *http.Request) {
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
