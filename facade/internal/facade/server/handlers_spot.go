package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (s *Server) getSpotPrice(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting spot price")

	resp, err := s.SpotService.GetSpotPrice(r)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getAllSpots(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all spots")

	resp, err := s.SpotService.GetAll(r)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getSpotByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting spot by ID")

	vars := mux.Vars(r)
	requestedSpotID := vars["id"]

	resp, err := s.SpotService.GetSpot(requestedSpotID)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) deleteSpotByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting spot by ID")

	vars := mux.Vars(r)
	requestedSpotID := vars["id"]

	resp, err := s.SpotService.DeleteSpot(requestedSpotID)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) addSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding spot")

	resp, err := s.SpotService.AddSpot(r)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) editSpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing spot")

	resp, err := s.SpotService.EditSpot(r)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

type availabilityInput struct {
	SpotIDs   []string  `json:"spot_ids" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

func (s *Server) getAvailableSpots(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting available spots")
	var input availabilityInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	if err := s.Validator.Struct(input); err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	checkIfExistInput := map[string]interface{}{
		"spot_ids": input.SpotIDs,
	}
	spotBody, err := json.Marshal(checkIfExistInput)
	if err != nil {
		s.handleError(w, "Failed to encode request body", err, http.StatusInternalServerError)
		return
	}

	spotResp, err := s.SpotService.CheckIfSpotsExist(spotBody)
	if err != nil {
		s.handleError(w, "Failed to send request to spot service", err, http.StatusInternalServerError)
		return
	}
	defer spotResp.Body.Close()

	spotBytes, err := io.ReadAll(spotResp.Body)
	if err != nil {
		s.handleError(w, "Failed to read response body", err, http.StatusInternalServerError)
		return
	}

	var result struct {
		NotFound []string `json:"not_found"`
		AllExist bool     `json:"all_exist"`
	}
	if err := json.Unmarshal(spotBytes, &result); err != nil {
		s.handleError(w, "Failed to parse response body", err, http.StatusInternalServerError)
		return
	}

	if !result.AllExist {
		s.handleError(w, "Some spots do not exist: "+strings.Join(result.NotFound, ", "), nil, http.StatusBadRequest)
		return
	}

	reservationBody, err := json.Marshal(input)
	if err != nil {
		s.handleError(w, "Failed to encode request body", err, http.StatusInternalServerError)
		return
	}

	// Forward the request to the reservation service
	resp, err := s.ReservationService.CheckAvailability(reservationBody)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	// Forward the response back to the user
	s.forwardResponse(w, resp)
}
