package http

import (
	"log"
	"net/http"

	gohtmx "github.com/falagansoftware/go-htmx/internal"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	router *mux.Router
	// Services used in routes
	UserService gohtmx.UserService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}

	// Handle if Panic
	s.router.Use(reportPanic)
	// Log Request
	s.router.Use(s.logRequest)

	// Register Router for unauthenticated routes
	r := s.router.PathPrefix("/").Subrouter()
	r.Use(s.requiredNoAuth)
	s.registerUserRoutes(r)
	return s
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
func reportPanic(next http.Handler) http.Handler {
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

// Middleware to log requests
func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Printf("[Request] %s %s %s %s", r.RemoteAddr, r.Method, r.URL, r.Proto)
		next.ServeHTTP(w, r)
	})
}
