package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
)

type Server struct {
	Router *chi.Mux
}

func New() *Server {
	s := &Server{
		Router: NewRouter(),
	}

	s.Router.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello World")) })
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
