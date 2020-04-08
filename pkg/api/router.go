package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		// Secure,
		middleware.StripSlashes,
		middleware.Compress(5),
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		Cors(),
	)

	return router
}
