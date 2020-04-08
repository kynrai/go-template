package server

import (
	"net/http"
	"regexp"

	"github.com/go-chi/cors"
)

var allowedOrigins = []*regexp.Regexp{
	regexp.MustCompile(`^(fluit\.co|localhost)(:(\d+))?$`),
	regexp.MustCompile(`.*\.(fluit\.co|localhost)(:(\d+))?$`),
}

func originValidator(r *http.Request, origin string) bool {
	for _, r := range allowedOrigins {
		if r.MatchString(origin) {
			return true
		}
	}
	return false
}

func Cors() func(http.Handler) http.Handler {
	// Common middleware
	return cors.New(cors.Options{
		AllowOriginFunc:  originValidator,
		AllowedMethods:   []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler
}
