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
		msg := "Failed to decode request body"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
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
		msg := "Failed to add user to MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
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
		msg := "Failed to decode request body"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditUser(input)
	if err != nil {
		msg := "Failed to edit user in MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
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
		msg := "Missing user ID"
		s.Logger.Error.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			msg := "User not found"
			s.Logger.Error.Printf("%s: %v", msg, id)
			http.Error(w, msg, http.StatusNotFound)
			return
		}

		msg := "Failed to delete user"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
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
		msg := "Missing user ID"
		s.Logger.Error.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	user, err := s.MongoDB.GetUser(id)
	if err != nil {
		msg := "Failed to get user"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		msg := "Failed to encode user to JSON"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User found: %v", user.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all users")
	users, err := s.MongoDB.GetAll()
	if err != nil {
		msg := "Failed to get users"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(users)
	if err != nil {
		msg := "Failed to encode users to JSON"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Users found: %v", len(users))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
