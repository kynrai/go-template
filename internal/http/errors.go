package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrBodyNotJSON = errors.New("body not valid json")

// compile time check to ensure HTTPError always implements StatusError
var _ StatusError = (*HTTPError)(nil)

// StatusError composes error and a method to return a status code
type StatusError interface {
	error
	Status() int
	JSON() string
}

// Error details
// swagger:response errorResp
type ErrorResp struct {
	// in:body
	Payload ErrorResponse
}

// HTTPError implements StatusError, returns an error and a HTTP status code
type HTTPError struct {
	Code int
	Err  error
}

// Error implements the error interface for HTTPError
func (s HTTPError) Error() string {
	if s.Err != nil {
		return s.Err.Error()
	}

	if status := s.Status(); status >= 400 {
		return http.StatusText(s.Status())
	}

	return ""
}

// Status implements the StatusError interface for HPPTError
func (s HTTPError) Status() int {
	return s.Code
}

type (
	ErrorField struct {
		Detail string `json:"detail"`
	}
	ErrorResponse struct {
		Errors []ErrorField `json:"errors"`
	}
)

func (s HTTPError) JSON() string {
	e := ErrorResponse{
		Errors: []ErrorField{
			{
				Detail: s.Error(),
			},
		},
	}
	j, _ := json.Marshal(e)
	return string(j)
}

// ErrorHandler returns an error from a http handler
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

// ServerHTTP implements the http.Handler interface, checks for an error and parses it if it is a StatusError
func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// response is always application/json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := h(w, r); err != nil {
		switch e := err.(type) {
		case StatusError:
			// We can retrieve the status here and write out a specific HTTP status code.
			w.WriteHeader(e.Status())
			fmt.Fprint(w, e.JSON())
		default:
			status := http.StatusInternalServerError
			if err == ErrBodyNotJSON {
				status = http.StatusBadRequest
			}
			w.WriteHeader(status)
			fmt.Fprint(w, (HTTPError{Code: status}).JSON())
		}
	}
}
