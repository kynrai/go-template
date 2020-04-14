package auth

import "net/http"

type Validator interface {
	Validate(h http.Handler) http.Handler
}
