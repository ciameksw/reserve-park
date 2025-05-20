package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type registerInput struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Registering a new user")
	var input registerInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	if err := s.Validator.Struct(input); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	input.Role = "user"

	modifiedBody, err := json.Marshal(input)
	if err != nil {
		s.handleError(w, "Failed to encode modified request body", err, http.StatusInternalServerError)
		return
	}

	resp, err := s.UserService.Register(modifiedBody)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Login user")

	resp, err := s.UserService.Login(r)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

type editInput struct {
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty"`
}

func (s *Server) editUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing a user")
	var input editInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	if err := s.Validator.Struct(input); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != input.UserID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	validatedBody, err := json.Marshal(input)
	if err != nil {
		s.handleError(w, "Failed to encode validated request body", err, http.StatusInternalServerError)
		return
	}

	resp, err := s.UserService.Edit(validatedBody)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

type editRoleInput struct {
	UserID string `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required,oneof=admin user"`
}

func (s *Server) editUsersRole(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing a user")
	var input editRoleInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	if err := s.Validator.Struct(input); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	validatedBody, err := json.Marshal(input)
	if err != nil {
		s.handleError(w, "Failed to encode validated request body", err, http.StatusInternalServerError)
		return
	}

	resp, err := s.UserService.Edit(validatedBody)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all users")

	resp, err := s.UserService.GetAll(r)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting user by userID")

	vars := mux.Vars(r)
	requestedUserID := vars["id"]

	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != requestedUserID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.UserService.GetUser(requestedUserID)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) deleteUserByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting user by userID")

	vars := mux.Vars(r)
	requestedUserID := vars["id"]

	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != requestedUserID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.UserService.DeleteUser(requestedUserID)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}
