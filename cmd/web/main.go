package main

import (
	"log"
	"net/http"
)

func main() {

	// create a new servemux and register the home function to act as the handler
	// for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/company", showCompany)
	mux.HandleFunc("/company/request", requestCompany)

	// Create a file server to serve the static frontend files and register the route
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
