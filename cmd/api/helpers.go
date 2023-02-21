package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// alphaNumeric will check the URL query paramter to ensure only alphanumberic characters are present
func alphaNumeric(name string) (string, bool) {
	return name, regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(name)
}

// readNameParam will grab the parameter from the URL using the request context
func (app *application) readNameParam(r *http.Request) (string, error) {
	// retrieve a slice containing any interpolated parameter names and values
	params := httprouter.ParamsFromContext(r.Context())

	// retrieve the name parameter
	name, ok := alphaNumeric(params.ByName("name"))
	if !ok || name == "" {
		return "", errors.New("invalid name parameter")
	}

	return name, nil
}

// envelope will help to add a layer to the JSON object for additional security measures
type envelope map[string]any

// writeJSON is a helper function which will iterate through the data passed and convert it into
// a JSON object to return.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append a newline to make the JSON more viewer-friendly
	js = append(js, '\n')

	// iterate through map and write the Header for each key-value pair
	for key, value := range headers {
		w.Header()[key] = value
	}

	// sets a response header so client knows response contains json
	w.Header().Set("Content-Type", "application/json")

	// write the status code for the response
	w.WriteHeader(status)
	// write the JSON as the HTTP response body
	w.Write(js)

	return nil
}

// readJSON will read the request body into a struct and check for different types of errors
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination
	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must be no larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will // return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message. err = dec.Decode(&struct{}{})
	err := dec.Decode(struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
