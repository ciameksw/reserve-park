package server

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Forwarding login request to user service")

	resp, err := s.UserService.Login(r)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Forwarding login request to user service")

	resp, err := s.UserService.GetAll(r)
	if err != nil {
		s.handleError(w, "Failed to send request to user service", err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (s *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
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
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (s *Server) deleteUserByID(w http.ResponseWriter, r *http.Request) {
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
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// Helper function to handle errors
func (s *Server) handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		s.Logger.Error.Printf("%s: %v", message, err)
	} else {
		s.Logger.Error.Println(message)
	}
	http.Error(w, message, statusCode)
}
