package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kai-munekuni/user-api/internal/domain/repository"
)

// Server struct of http server
type Server struct {
	server *http.Server
}

// NewServer initialize server
func NewServer(port int, ur repository.User) *Server {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/ping", &healthCheckHandler{}).Methods("GET")
	r.Handle("/signup", &signupHandler{userRepo: ur}).Methods("POST")

	authedRouter := r.PathPrefix("/").Subrouter()
	m := authMiddleware{userRepo: ur}
	authedRouter.Use(m.Middleware)
	authedRouter.Handle("/users/{id}", &getUserHandler{userRepo: ur}).Methods("GET")
	authedRouter.Handle("/users/{id}", &patchUserHandler{userRepo: ur}).Methods("PATCH")
	authedRouter.Handle("/close", &deleteUserHandler{userRepo: ur}).Methods("POST")

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}
}

// Start start http server
func (s *Server) Start() error {
	log.Println("Server running...")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop stop http server
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutdown server")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}

	return nil
}
