package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/kynrai/go-template/pkg/auth"
)

type Server struct {
	Router *chi.Mux
	Auth   auth.Auth0Repo
}

func New() *Server {
	s := &Server{
		Router: NewRouter(),
		Auth:   auth.NewAuth0(), // Auth 0 instance
		// Auth:   auth.NewFirebase(), // Firebase instance
	}

	s.Router.Route("/auth0-protected", func(r chi.Router) {
		r.Use(s.Auth.Validate)
	})

	s.Router.Route("/v1", func(r chi.Router) {
	})

	return s
}

func (s *Server) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not set, defaulting to 5000")
		port = "5000"
	}
	log.Printf("serving on port %s\n", port)
	server := http.Server{Addr: ":" + port, Handler: s.Router}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	server.Shutdown(context.Background())
}
