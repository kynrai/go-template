package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	auth "github.com/kynrai/go-template/pkg/auth/auth0"
	auth0 "github.com/kynrai/go-template/pkg/auth/auth0"
	firebaseauth "github.com/kynrai/go-template/pkg/auth/firebase"
)

type Server struct {
	Router       *chi.Mux
	Auth0        auth.Repo
	FirebaseAuth firebaseauth.Repo
}

func New() *Server {
	s := &Server{
		Router:       NewRouter(),
		Auth0:        auth0.New(),            // Auth 0 instance
		FirebaseAuth: firebaseauth.MustNew(), // Firebase instance
	}

	s.Router.Route("/auth0-protected", func(r chi.Router) {
		r.Use(s.Auth0.Validate)
	})

	s.Router.Route("/firebase-protected", func(r chi.Router) {
		r.Use(firebaseauth.Validate(s.FirebaseAuth))
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
