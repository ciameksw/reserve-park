package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ciameksw/reserve-park/user/internal/user/auth"
	m "github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type addInput struct {
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Role     m.RoleType `json:"role"`
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding user")
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	existingUser, err := s.MongoDB.GetUserByUsernameOrEmail(input.Username, input.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		s.handleError(w, "Failed to check for existing user", err, http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		s.handleError(w, "Username or email already exists", nil, http.StatusConflict)
		return
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		s.handleError(w, "Failed to hash the password", err, http.StatusInternalServerError)
		return
	}

	data := m.User{
		UserID:       uuid.NewString(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         input.Role,
		UpdatedAt:    time.Now(),
	}

	if err := s.Validator.Struct(data); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.AddUser(data)
	if err != nil {
		s.handleError(w, "Failed to add user to MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User added: %v", data.Username)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(data.UserID))
}

type editInput struct {
	UserID   string     `json:"user_id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Role     m.RoleType `json:"role"`
}

func (s *Server) editUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing user")
	var input editInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	// TODO: Change this, we dont want to provide the password every time
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		s.handleError(w, "Failed to hash the password", err, http.StatusInternalServerError)
		return
	}

	data := m.User{
		UserID:       uuid.NewString(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         input.Role,
		UpdatedAt:    time.Now(),
	}

	if err := s.Validator.Struct(data); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditUser(data)
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

type loginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Logging user")
	var input loginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	if err := s.Validator.Struct(input); err != nil {
		s.handleError(w, "Invalid input data", err, http.StatusBadRequest)
		return
	}

	user, err := s.MongoDB.GetUserByUsernameOrEmail(input.Username, "")
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Invalid username or password", err, http.StatusUnauthorized)
			return
		}

		s.handleError(w, "Unexpected server error", err, http.StatusInternalServerError)
		return
	}

	match := auth.VerifyPassword(input.Password, user.PasswordHash)
	if !match {
		s.handleError(w, "Invalid username or password", err, http.StatusUnauthorized)
		return
	}

	jwt, err := auth.GenerateJWT(user.UserID, user.Role, s.Config.Salt)
	if err != nil {
		s.handleError(w, "Failed to generate JWT", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User logged in: %v", user.Username)
	resp := map[string]string{
		"jwt": jwt,
	}
	s.writeJSON(w, resp, http.StatusOK)
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
