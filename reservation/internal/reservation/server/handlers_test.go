package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ciameksw/reserve-park/reservation/internal/reservation/config"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/logger"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/mongodb"
	"github.com/gorilla/mux"
)

var s *Server
var reservationID string
var userID = "75390349821"
var spotID = "96363829890"

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

func TestAddReservation(t *testing.T) {
	input := addInput{
		UserID:    userID,
		SpotID:    spotID,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Status:    "valid",
		PricePaid: 10.0,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/reservations", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.addReservation)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	reservationID = rr.Body.String()
}

func TestGetAllReservations(t *testing.T) {
	req, err := http.NewRequest("GET", "/reservations", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getAllReservations)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var reservations []mongodb.Reservation
	err = json.NewDecoder(rr.Body).Decode(&reservations)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(reservations) != 1 {
		t.Errorf("handler returned wrong number of reservations: got %v want %v", len(reservations), 1)
	}
}

func TestGetReservation(t *testing.T) {
	req, err := http.NewRequest("GET", "/reservations/"+reservationID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/reservations/{id}", s.getReservation).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestEditReservation(t *testing.T) {
	input := editInput{
		ReservationID: reservationID,
		UserID:        userID,
		SpotID:        spotID,
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(2 * time.Hour),
		Status:        "valid",
		PricePaid:     10.0,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("PUT", "/reservations", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.editReservation)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestCheckAvailabilityOccupied(t *testing.T) {
	input := mongodb.AvailabilityInput{
		SpotIDs:   []string{spotID},
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("GET", "/reservations/availability/check", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.checkAvailability)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var availableSpots []string
	err = json.NewDecoder(rr.Body).Decode(&availableSpots)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(availableSpots) != 0 {
		t.Errorf("handler returned wrong number of available spots: got %v want %v", len(availableSpots), 0)
	}
}

func TestCheckAvailabilityFree(t *testing.T) {
	input := mongodb.AvailabilityInput{
		SpotIDs:   []string{spotID},
		StartTime: time.Now().Add(2 * time.Hour),
		EndTime:   time.Now().Add(3 * time.Hour),
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("GET", "/reservations/availability/check", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.checkAvailability)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var availableSpots []string
	err = json.NewDecoder(rr.Body).Decode(&availableSpots)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(availableSpots) != 1 {
		t.Errorf("handler returned wrong number of available spots: got %v want %v", len(availableSpots), 1)
	}
}

func TestGetReservationsByUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/reservations/user/"+userID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/reservations/user/{id}", s.getUserReservations).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var reservations []mongodb.Reservation
	err = json.NewDecoder(rr.Body).Decode(&reservations)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(reservations) != 1 {
		t.Errorf("handler returned wrong number of reservations: got %v want %v", len(reservations), 1)
	}
}

func TestGetReservationsBySpot(t *testing.T) {
	req, err := http.NewRequest("GET", "/reservations/spot/"+spotID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/reservations/spot/{id}", s.getSpotReservations).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var reservations []mongodb.Reservation
	err = json.NewDecoder(rr.Body).Decode(&reservations)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(reservations) != 1 {
		t.Errorf("handler returned wrong number of reservations: got %v want %v", len(reservations), 1)
	}
}

func TestDeleteReservation(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/reservations/"+reservationID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/reservations/{id}", s.deleteReservation).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
