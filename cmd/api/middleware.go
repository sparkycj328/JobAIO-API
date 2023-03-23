package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type any, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// rateLimit will use golang's x/time/rate to declare a new limiter and then use
// this limiter upon each request handled by the http client
func (app *application) rateLimit(next http.Handler) http.Handler {
	// struct that serves the same purpose as the old clients map
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// launch a go routine that will remove old entries from the clients IP map once every minute
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent concurrent access to the clients map
			mu.Lock()

			// loop through the client IPs and remove any entries for clients
			// that have not been seen within the last three minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			// unlock the mutex once cleanup is complete
			mu.Unlock()
		}
	}()
	// return an http handler which wraps around each request through the http client
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only carry out the check if rate limiting is enabled
		if app.config.limiter.enabled {
			// grab the ip address completing the request
			// return a server error if unable to read the IP address
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Lock the mutex to prevent concurrent writing to the clients map
			mu.Lock()

			// Check if the ip address already exists in the map
			// if it does not exist, add a new limiter to the clients map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
			}

			// Update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// if the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
			// response, mutex will be unlocked before sending the 429 response
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceeded(w, r)
				return
			}

			// unlock the mutex before calling the next handler in the chain
			mu.Unlock()

		}
		next.ServeHTTP(w, r)
	})

}
