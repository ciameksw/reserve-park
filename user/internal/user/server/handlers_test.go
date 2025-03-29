package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ciameksw/reserve-park/user/internal/user/config"
	"github.com/ciameksw/reserve-park/user/internal/user/logger"
	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/gorilla/mux"
)

var s *Server
var userID string
var jwt string
var newUsername = "NewUsername"
var newPassword = "NewPassword"

func TestMain(m *testing.M) {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	// Connect to mock MongoDB
	db, err := mongodb.ConnectMock()
	if err != nil {
		lgr.Error.Fatalf("Failed to connect to mock MongoDB: %v", err)
	}
	defer db.Disconnect()

	s = NewServer(lgr, cfg, db)

	os.Exit(m.Run())
}

func TestAddUser(t *testing.T) {
	input := addInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     "user",
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.addUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	userID = rr.Body.String()
}

func TestGetAllUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getAllUsers)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []mongodb.User
	err = json.NewDecoder(rr.Body).Decode(&users)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("handler returned wrong number of users: got %v want %v", len(users), 1)
	}
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/"+userID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", s.getUser).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestEditUser(t *testing.T) {
	input := editInput{ // Partial edit
		UserID:   userID,
		Username: newUsername,
		Email:    "test123@example.com",
		Password: newPassword,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("PATCH", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.editUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLogin(t *testing.T) {
	input := loginInput{
		Username: newUsername,
		Password: newPassword,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var output loginResponse
	err = json.NewDecoder(rr.Body).Decode(&output)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	jwt = output.Jwt
	fmt.Println(jwt)
}

func TestAuthorize(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/authorize", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.authorize)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check the response content
	if response["role"] != "user" {
		t.Errorf("handler returned wrong response: got %v want %v", response["user"], true)
	}
}

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/"+userID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", s.deleteUser).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
