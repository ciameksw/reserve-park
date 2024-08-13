package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ciameksw/reserve-park/user/internal/user/config"
	"github.com/ciameksw/reserve-park/user/internal/user/logger"
	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	m "github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/gorilla/mux"
)

var s *Server
var userID string

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

	var interfaceSlice []interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &interfaceSlice)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(interfaceSlice) != 1 {
		t.Errorf("handler returned wrong number of users: got %v want %v", len(interfaceSlice), 1)
	}

	userID = interfaceSlice[0].(map[string]interface{})["user_id"].(string)
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
	input := m.User{
		UserID:   userID,
		Username: "editeduser",
		Email:    "test@example.com",
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("PUT", "/users", bytes.NewBuffer(body))
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

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/"+userID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", s.getUser).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
