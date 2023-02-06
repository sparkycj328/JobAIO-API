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

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
