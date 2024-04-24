package http

import (
	"log"
	"net/http"
	"strconv"

	gohtmx "github.com/falagansoftware/go-htmx/internal"
	html "github.com/falagansoftware/go-htmx/internal/http/html"
)

// Routes

func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/users", s.handleUserList).Methods("GET")
	s.router.HandleFunc("/users/filter", s.handleUserListFilter).Methods("GET")
}

// Handlers

func (s *Server) handleUserList(w http.ResponseWriter, r *http.Request) {
	// Filter
	filter := getUserFilters(r)
	// Get all users
	users, err := s.UserService.FindUsers(r.Context(), filter)
	log.Print(filter)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
	// Render Users
	view := html.UserList(users, filter.Sort, filter.Order)
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

func (s *Server) handleUserListFilter(w http.ResponseWriter, r *http.Request) {
	// Filter
	filter := getUserFilters(r)
	// Get all users
	users, err := s.UserService.FindUsers(r.Context(), filter)
	log.Print(filter)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
	// Render Users
	view := html.UserListSync(users, filter.Sort, filter.Order)
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

// Helpers

func getUserFilters(r *http.Request) *gohtmx.UserFilters {
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
	}

	sort := r.URL.Query().Get("sort")

	if sort != "" {
		filter.Sort = sort
	} else {
		filter.Sort = "name"
	}

	order := r.URL.Query().Get("order")

	if order != "" {
		filter.Order = order
	} else {
		filter.Order = "ASC"
	}

	return &filter
}
