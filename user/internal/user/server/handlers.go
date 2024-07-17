package server

import (
	"encoding/json"
	"fmt"
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
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) editUser(w http.ResponseWriter, r *http.Request) {
	var input m.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditUser(input)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	user, err := s.MongoDB.GetUser(id)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.MongoDB.GetAll()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
