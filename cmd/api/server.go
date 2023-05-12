package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve
func (app *application) serve() error {

	// Initialize the local server with our desired configuration
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Use a shutdown error channel to receive any errors returned by the shutdown function.
	shutdownError := make(chan error)

	// Anonymous function will listen in background for certain signals to allow for graceful shutdown
	go func() {
		// Create a channel to receive SIGINT and SIGTERM signals
		// initial channel must be buffered in order to receive signals
		quit := make(chan os.Signal, 1)

		// Listen for incoming SIGINT and SIGTERM calls and relay them
		// to the quit channel. No other signals will be relayed
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// block until a signal is received
		s := <-quit

		// If a signal is received, print a message to the console
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// Create a context with a 20-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// call the shutdown function, passing it the context to initiate graceful shutdown
		// if the error returns nil then graceful shutdown was successful, otherwise
		// the server had issues closing open connections, or it exceeded the 20-second timeout
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		app.wg.Wait()
		shutdownError <- nil
	}()

	// Log a startup message and start the server as normal
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Start the local server and return an error
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. will only return an
	// error if it is NOT http.ErrServerClosed()
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise receive the value from the shutdownError on the shutdownError channel
	if err = <-shutdownError; err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
