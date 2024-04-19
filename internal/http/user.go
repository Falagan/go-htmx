package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	gohtmx "github.com/falagansoftware/go-htmx/internal"
	"github.com/gorilla/mux"
)

// Routes

func (s *Server) registerUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", s.handleUserList).Methods("GET")
}

// Handlers

func (s *Server) handleUserList(w http.ResponseWriter, r *http.Request) {
	// Filters
	filter := gohtmx.UserFilters{}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 20
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

	active,_ := strconv.ParseBool(r.URL.Query().Get("active"))

	if active {
		filter.Active = active
	}

	// Get all users
	users, err := s.UserService.FindUsers(r.Context(), filter)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the users
	err = json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
