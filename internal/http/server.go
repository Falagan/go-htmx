package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	gohtmx "github.com/falagansoftware/go-htmx/internal"
	"github.com/gorilla/mux"
)

type Server struct {
	server  *http.Server
	router  *mux.Router
	Address string
	Port    int
	// Services used in routes
	UserService gohtmx.UserService
}

type Option func(*Server)

func NewServer(address string, port int, options ...Option) *Server {
	s := &Server{
		server:  &http.Server{},
		router:  mux.NewRouter(),
		Address: address,
		Port:    port,
	}

	for _, option := range options {
		option(s)
	}

	// Handle if Panic
	s.router.Use(s.reportPanic)
	// Log Request
	s.router.Use(s.logRequest)
	// Attach router
	// s.server.Handler = s.router

	// Register Router for unauthenticated routes
	// r := s.router.PathPrefix("/").Subrouter()
	// r.Use(s.requiredNoAuth)
	s.registerUserRoutes()
	return s
}

func (s *Server) ListenAndServe() error {
	// h.Handle("/metrics", promhttp.Handler())
	url := s.URL()
	return http.ListenAndServe(url, s.router)
}

// Options

func WithTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.server.WriteTimeout = timeout
	}
}

// Middlewares

// Auth

// Middleware to check that user is not authenticated
func (s *Server) requiredNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a valid session
		// If session is valid, redirect to /dashboard
		// If session is invalid, call next.ServeHTTP(w, r)
	})
}

// Middleware to check that user is authenticated
func (s *Server) requiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a valid session
		// If session is valid, call next.ServeHTTP(w, r)
		// If session is invalid, redirect to /login
	})
}

// Panic
func (s *Server) reportPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[Panic] %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Logs

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Printf("[Request] %s %s %s %s", r.RemoteAddr, r.Method, r.URL, r.Proto)
		next.ServeHTTP(w, r)
	})
}

// Helpers

func (s *Server) URL() string {
	return s.Address + ":" + strconv.Itoa(s.Port)
}
