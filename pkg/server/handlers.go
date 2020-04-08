package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	h "github.com/kynrai/go-template/internal/http"
)

func Health(msgs ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for _, msg := range msgs {
			w.Write([]byte(msg))
		}
	}
}

func BodyCloser(hl http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		hl.ServeHTTP(w, r)
	})
}

// PanicHandler intercepts a panic in the handler and prints out debugging information
func PanicHandler(hl http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				headerStr := make([]string, 0, len(r.Header))
				for k, v := range r.Header {
					headerStr = append(headerStr, fmt.Sprintf("%s: %s", k, v))
				}
				log.Printf("Panic in handler %s %s: %+v\n", r.Method, r.URL, rec)
				log.Printf("%s\n", string(debug.Stack()))

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(h.HTTPError{Code: http.StatusInternalServerError, Err: errors.New("Internal error")}.JSON()))
			}
		}()
		hl.ServeHTTP(w, r)
	})
}
