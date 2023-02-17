package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// config struct will hold all configuration settings for our applications
type config struct {
	port int
	env  string
}

// application struct will hold the dependencies for our HTTP handlers
// helper functions and middleware
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// read flag values into config struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)\"")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// declares an instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	// create a new servemux and register the healthcheck function
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, cfg.port)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
