package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ciameksw/reserve-park/spot/internal/spot/config"
	"github.com/ciameksw/reserve-park/spot/internal/spot/logger"
	"github.com/ciameksw/reserve-park/spot/internal/spot/mongodb"
	"github.com/gorilla/mux"
)

var s *Server
var spotID string
var pricePerHour float64

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
		Size:         mongodb.SizeLarge,
		Type:         mongodb.SpotTypeOutdoor,
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

	spotID = rr.Body.String()
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

	var spots []mongodb.Spot
	err = json.NewDecoder(rr.Body).Decode(&spots)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(spots) != 1 {
		t.Errorf("handler returned wrong number of spots: got %v want %v", len(spots), 1)
	}
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
	pricePerHour = 10.5
	latitude := -34.7365
	longitude := 86.8271
	input := editInput{
		SpotID:       spotID,
		Latitude:     &latitude,
		Longitude:    &longitude,
		PricePerHour: &pricePerHour,
		Size:         mongodb.SizeSmall,
		Type:         mongodb.SpotTypeEV,
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("PUT", "/spots", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.editSpot)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestGetPrice(t *testing.T) {
	input := mongodb.GetPriceInput{
		SpotID:    spotID,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
	}
	body, _ := json.Marshal(input)
	req, err := http.NewRequest("GET", "/spots/price", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.getPrice)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["spot_id"] != spotID {
		t.Errorf("handler returned wrong spotID: got %v want %v", response["spot_id"], spotID)
	}

	correctPrice := 2 * pricePerHour
	if response["price"] != correctPrice {
		t.Errorf("handler returned wrong price: got %v want %v", response["price"], correctPrice)
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
