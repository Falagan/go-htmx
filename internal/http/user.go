package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	gohtmx "github.com/falagansoftware/go-htmx/internal"
)

// Routes

func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/users", s.handleUserList).Methods("GET")
}

// Handlers

func (s *Server) handleUserList(w http.ResponseWriter, r *http.Request) {
	// Filters
	filter := gohtmx.UserFilters{}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		filter.Limit = 20
	} else {
		filter.Limit = limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

	if err != nil {
		offset = 0
	} else {
		filter.Offset = offset
	}

	name := r.URL.Query().Get("name")

	if name != "" {
		filter.Name = &name
	}

	surname := r.URL.Query().Get("surname")

	if surname != "" {
		filter.Surname = &surname
	}

	email := r.URL.Query().Get("email")

	if email != "" {
		filter.Email = &email
	}

	active := r.URL.Query().Get("active")

	if active != "" {
		filter.Active, _ = strconv.ParseBool(active)
	} else {
		filter.Active = true
	}

	// Get all users
	users, err := s.UserService.FindUsers(r.Context(), filter)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}

	// Render the users
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
	}
	
}
