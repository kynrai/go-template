package auth

import (
	"log"
	"net/http"
	"strings"
)

func Validate(v Validator) func(h http.Handler) http.Handler {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			headerParts := strings.Split(r.Header.Get("Authorization"), " ")
			// we expect a Bearer Token
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token, err := v.Validate(r.Context(), headerParts[1])
			if err != nil {
				log.Println("Failed to verify token", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			// all good so add token and serve handler
			r = r.WithContext(NewContextWithToken(r.Context(), token))
			f.ServeHTTP(w, r)
		})
	}
}
