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
	m "github.com/ciameksw/reserve-park/reservation/internal/reservation/mongodb"
	"github.com/gorilla/mux"
)

var s *Server
var reservationID string

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
		UserID: "12345",
		SpotID: "12345",
		Start:  time.Now(),
		End:    time.Now().Add(time.Hour),
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

	var reservations []m.Reservation
	err = json.NewDecoder(rr.Body).Decode(&reservations)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(reservations) != 1 {
		t.Errorf("handler returned wrong number of reservations: got %v want %v", len(reservations), 1)
	}

	reservationID = reservations[0].ReservationID
}

func TestGetUser(t *testing.T) {
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
	input := m.Reservation{
		ReservationID: reservationID,
		UserID:        "54321",
		SpotID:        "54321",
		Start:         time.Now(),
		End:           time.Now().Add(time.Hour),
		Canceled:      true,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("PUT", "/reservations", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.editReservation)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
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
