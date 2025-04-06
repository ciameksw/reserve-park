package server

import (
	"net/http"

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
