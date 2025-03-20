package server

import (
	"encoding/json"
	"net/http"

	m "github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type addInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding user")
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	data := m.User{
		UserID:          uuid.NewString(),
		Username:        input.Username,
		Email:           input.Email,
		TotalMoneySpent: 0,
	}

	err = s.MongoDB.AddUser(data)
	if err != nil {
		s.handleError(w, "Failed to add user to MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User added: %v", data.Username)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) editUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing user")
	var input m.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditUser(input)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "User not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to edit user in MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User edited: %v", input.Username)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting user")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing user ID", nil, http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "User not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to delete user", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User deleted: %v", id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting user")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing user ID", nil, http.StatusBadRequest)
		return
	}

	user, err := s.MongoDB.GetUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "User not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User found: %v", user.Username)
	s.writeJSON(w, user, http.StatusOK)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all users")
	users, err := s.MongoDB.GetAll()
	if err != nil {
		s.handleError(w, "Failed to get users", err, http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		s.Logger.Info.Println("No users found")
		s.writeJSON(w, []m.User{}, http.StatusOK)
		return
	}

	s.Logger.Info.Printf("Users found: %v", len(users))
	s.writeJSON(w, users, http.StatusOK)
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

// Helper function to write JSON responses
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	j, err := json.Marshal(data)
	if err != nil {
		s.handleError(w, "Failed to encode response to JSON", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(j)
}
