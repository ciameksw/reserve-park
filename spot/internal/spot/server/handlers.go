package server

import (
	"encoding/json"
	"net/http"

	m "github.com/ciameksw/reserve-park/spot/internal/spot/mongodb"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type addInput struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	PricePerHour float64 `json:"price_per_hour"`
}

func (s *Server) addSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding spot")
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := "Failed to decode request body"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	data := m.Spot{
		SpotID:       uuid.NewString(),
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		PricePerHour: input.PricePerHour,
	}

	err = s.MongoDB.AddSpot(data)
	if err != nil {
		msg := "Failed to add spot to MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot added: %v", data.SpotID)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) editSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing spot")
	var input m.Spot

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := "Failed to decode request body"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditSpot(input)
	if err != nil {
		msg := "Failed to edit spot in MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
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
		msg := "Spot ID not provided"
		s.Logger.Error.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteSpot(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			msg := "Spot not found"
			s.Logger.Error.Println(msg)
			http.Error(w, msg, http.StatusNotFound)
			return
		}

		msg := "Failed to delete spot"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot deleted: %v", id)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting spot")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		msg := "Spot ID not provided"
		s.Logger.Error.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	spot, err := s.MongoDB.GetSpot(id)
	if err != nil {
		msg := "Failed to get spot"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(spot)
	if err != nil {
		msg := "Failed to encode spot to JSON"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spot found: %v", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *Server) getAllSpots(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all spots")

	spots, err := s.MongoDB.GetAll()
	if err != nil {
		msg := "Failed to get all spots"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(spots)
	if err != nil {
		msg := "Failed to encode spots to JSON"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Spots found: %v", len(spots))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
