package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// serve writes the status header and marshals the struct passed to include in the response body.
func serve(res http.ResponseWriter, val interface{}, code int) {
	if val != nil {
		b, err := json.Marshal(&val)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(code)
		res.Write(b)
		return
	}

	res.WriteHeader(code)
}

// statusError is returned when API calls fail.
type statusError struct {
	Err string `json:"error"`
}

// serveError writes the http header and error response.
// Note that when `err` is of type dbError or transError, `code` will be overridden.
func serveError(res http.ResponseWriter, err error, code int) {
	// We will marshal serr in the response.
	serr := statusError{
		Err: err.Error(),
	}

	// Set headers and encode the response.
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(code)
	if err := json.NewEncoder(res).Encode(serr); err != nil {
		log.Errorf("error encoding json error %s", err)
	}
}
