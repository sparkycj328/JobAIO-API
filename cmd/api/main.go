package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/sparkycj328/JobAIO-API/internal/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

// config struct will hold all configuration settings for our applications
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// application struct will hold the dependencies for our HTTP handlers
// helper functions and middleware
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	// read flag values into config struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)\"")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv(""), "Database connection")
	fmt.Println(&cfg.db.dsn)
	// read flag values to configure the database
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// call openDB and defer the db from closing until main finishes
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// declares an instance of the application struct
	// passes it our config, logger, and the database connection pool
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModel(db),
	}

	// declares an instance of an http Server where we can use the router
	// which can be passed upon declaration of the http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// openDB will open the designated db based on the dsn
// it will then ping the db and return it if no errors occurred
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish a new connection and if connection takes longer than 5 seconds
	// then the context will cancel and return an error
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
