package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ciameksw/reserve-park/spot/internal/spot/config"
	"github.com/ciameksw/reserve-park/spot/internal/spot/logger"
	"github.com/ciameksw/reserve-park/spot/internal/spot/mongodb"
	"github.com/gorilla/mux"
)

var s *Server
var spotID string

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

func TestAddSpot(t *testing.T) {
	input := addInput{
		Latitude:     34.7365,
		Longitude:    -86.8271,
		PricePerHour: 5.00,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/spots", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.addSpot)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetAllSpots(t *testing.T) {
	req, err := http.NewRequest("GET", "/spots", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getAllSpots)

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
		t.Errorf("handler returned wrong number of spots: got %v want %v", len(interfaceSlice), 1)
	}

	spotID = interfaceSlice[0].(map[string]interface{})["spot_id"].(string)
}

func TestGetSpot(t *testing.T) {
	req, err := http.NewRequest("GET", "/spots/"+spotID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/spots/{id}", s.getSpot).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestEditUser(t *testing.T) {
	input := mongodb.Spot{
		SpotID:    spotID,
		Latitude:  -34.7365,
		Longitude: 86.8271,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("PUT", "/spots", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.editSpot)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/spots/"+spotID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/spots/{id}", s.deleteSpot).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
