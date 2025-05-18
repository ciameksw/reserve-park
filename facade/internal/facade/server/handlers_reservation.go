package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) getAllReservations(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all reservations")

	resp, err := s.ReservationService.GetAll(r)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) deleteReservationByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting reservation by reservationID")

	vars := mux.Vars(r)
	requestedReservationID := vars["id"]

	resp, err := s.ReservationService.DeleteReservation(requestedReservationID)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getReservationsBySpot(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting reservations by spotID")

	vars := mux.Vars(r)
	spotID := vars["id"]

	resp, err := s.ReservationService.GetReservationsBySpot(spotID)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getReservationsByUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting reservations by userID")

	vars := mux.Vars(r)
	userID := vars["id"]

	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != userID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.ReservationService.GetReservationsByUser(userID)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) getReservationByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting reservation by reservationID")

	vars := mux.Vars(r)
	requestedReservationID := vars["id"]

	// Send the request to the reservation service
	resp, err := s.ReservationService.Get(requestedReservationID)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Perform authorization check early
	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	// Early return if the user is an admin
	if RoleType(authResp.Role) == RoleAdmin {
		s.Logger.Info.Println("Admin access granted")

		s.forwardResponse(w, resp)
		return
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.handleError(w, "Failed to read response body", err, http.StatusInternalServerError)
		return
	}

	// Parse the response body into a map
	var responseMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
		s.handleError(w, "Failed to parse response body", err, http.StatusInternalServerError)
		return
	}

	// Get the userID from the response
	userID, ok := responseMap["user_id"].(string)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	// Check if the user is authorized to access this reservation
	if authResp.UserID != userID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	// Replace the response body so it can be forwarded
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	s.forwardResponse(w, resp)
}

func (s *Server) addReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding reservation")

	// Read the request body into memory
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, "Failed to read request body", err, http.StatusBadRequest)
		return
	}

	// Parse the JSON body into a map to check the user_id
	var requestBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
		s.handleError(w, "Failed to parse request body", err, http.StatusBadRequest)
		return
	}

	// Extract the user_id from the request body
	userID, ok := requestBody["user_id"].(string)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusBadRequest)
		return
	}

	// Perform authorization check
	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != userID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	// Forward the request to the reservation service
	resp, err := s.ReservationService.Add(bodyBytes)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	// Forward the response back to the user
	s.forwardResponse(w, resp)
}

func (s *Server) editReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing reservation")

	// Read the request body into memory
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, "Failed to read request body", err, http.StatusBadRequest)
		return
	}

	// Parse the JSON body into a map to check the user_id
	var requestBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
		s.handleError(w, "Failed to parse request body", err, http.StatusBadRequest)
		return
	}

	// Extract the user_id from the request body
	userID, ok := requestBody["user_id"].(string)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusBadRequest)
		return
	}

	// Perform authorization check
	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != userID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	// Forward the request to the reservation service
	resp, err := s.ReservationService.Edit(bodyBytes)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	// Forward the response back to the user
	s.forwardResponse(w, resp)
}

func (s *Server) cancelReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Cancel reservation")

	vars := mux.Vars(r)
	requestedReservationID := vars["id"]

	// Send the request to the reservation service
	resp, err := s.ReservationService.Get(requestedReservationID)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.handleError(w, "Failed to read response body", err, http.StatusInternalServerError)
		return
	}

	// Parse the response body into a map
	var responseMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
		s.handleError(w, "Failed to parse response body", err, http.StatusInternalServerError)
		return
	}

	// Get the userID from the response
	userID, ok := responseMap["user_id"].(string)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	// Perform authorization check early
	authResp, ok := r.Context().Value(authorizeKey).(authorizeResponse)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if RoleType(authResp.Role) != RoleAdmin && authResp.UserID != userID {
		s.handleError(w, "Unauthorized", nil, http.StatusUnauthorized)
		return
	}

	// Get the status from the response
	status, ok := responseMap["status"].(string)
	if !ok {
		s.handleError(w, "Unexpected error", nil, http.StatusInternalServerError)
		return
	}

	if status == "canceled" {
		s.handleError(w, "Reservation already canceled", nil, http.StatusConflict)
		return
	}

	// Create a map with reservation ID and status
	cancelRequest := map[string]string{
		"reservation_id": requestedReservationID,
		"status":         "canceled",
	}
	cancelBytes, err := json.Marshal(cancelRequest)
	if err != nil {
		s.handleError(w, "Failed to encode cancel request body", err, http.StatusInternalServerError)
		return
	}

	// Send the request to the reservation service
	resp, err = s.ReservationService.Edit(cancelBytes)
	if err != nil {
		s.handleError(w, "Failed to send request to reservation service", err, http.StatusInternalServerError)
		return
	}

	// Forward the response back to the user
	s.forwardResponse(w, resp)
}
