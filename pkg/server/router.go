package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		BodyCloser,
		// Secure,
		middleware.StripSlashes,
		middleware.Compress(5),
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		PanicHandler,
		Cors(),
	)

	router.Method(http.MethodGet, "/health", Health("OK"))
	router.Method(http.MethodHead, "/health", Health())

	return router
}
