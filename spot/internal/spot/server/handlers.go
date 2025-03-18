package server

import (
	"encoding/json"
	"net/http"
	"time"

	m "github.com/ciameksw/reserve-park/spot/internal/spot/mongodb"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type addInput struct {
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	PricePerHour float64    `json:"price_per_hour"`
	Size         m.SizeType `json:"size"`
	Type         m.SpotType `json:"type"`
}

func (s *Server) addSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding spot")
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadGateway)
		return
	}

	data := m.Spot{
		SpotID:       uuid.NewString(),
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		PricePerHour: input.PricePerHour,
		Size:         input.Size,
		Type:         input.Type,
		UpdatedAt:    time.Now(),
	}

	if err := s.Validator.Struct(data); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.AddSpot(data)
	if err != nil {
		s.handleError(w, "Failed to add spot to MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot added: %v", data.SpotID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(data.SpotID))
}

type editInput struct {
	SpotID       string     `json:"spot_id"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	PricePerHour float64    `json:"price_per_hour"`
	Size         m.SizeType `json:"size"`
	Type         m.SpotType `json:"type"`
}

func (s *Server) editSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing spot")
	var input editInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	data := m.Spot{
		SpotID:       input.SpotID,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		PricePerHour: input.PricePerHour,
		Size:         input.Size,
		Type:         input.Type,
		UpdatedAt:    time.Now(),
	}

	if err := s.Validator.Struct(data); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditSpot(data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Spot not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to edit spot in MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot edited: %v", input.SpotID)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting spot")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing spot ID", nil, http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteSpot(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Spot not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to delete spot", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot deleted: %v", id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting spot")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing spot ID", nil, http.StatusBadRequest)
		return
	}

	spot, err := s.MongoDB.GetSpot(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Spot not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to get spot from MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot found: %v", id)
	s.writeJSON(w, spot, http.StatusOK)
}

func (s *Server) getAllSpots(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all spots")

	spots, err := s.MongoDB.GetAll()
	if err != nil {
		s.handleError(w, "Failed to get all spots", err, http.StatusInternalServerError)
		return
	}

	if len(spots) == 0 {
		s.Logger.Info.Println("No spots found")
		s.writeJSON(w, []m.Spot{}, http.StatusOK)
		return
	}

	s.Logger.Info.Printf("Spots found: %v", len(spots))
	s.writeJSON(w, spots, http.StatusOK)
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
