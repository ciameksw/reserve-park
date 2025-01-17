package server

import (
	"encoding/json"
	"net/http"
	"time"

	m "github.com/ciameksw/reserve-park/reservation/internal/reservation/mongodb"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type addInput struct {
	UserID string    `json:"user_id"`
	SpotID string    `json:"spot_id"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

func (s *Server) addReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding reservation")
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := "Failed to decode request body"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	data := m.Reservation{
		ReservationID: uuid.NewString(),
		UserID:        input.UserID,
		SpotID:        input.SpotID,
		Start:         input.Start,
		End:           input.End,
		Canceled:      false,
	}

	err = s.MongoDB.AddReservation(data)
	if err != nil {
		msg := "Failed to add reservation to MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Reservation added: %v", data.ReservationID)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) editReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Editing reservation")
	var input m.Reservation

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := "Failed to decode request body"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditReservation(input)
	if err != nil {
		msg := "Failed to edit reservation in MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Reservation edited: %v", input.ReservationID)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting reservation")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		msg := "Missing reservation ID"
		s.Logger.Error.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteReservation(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			msg := "Reservation not found"
			s.Logger.Error.Printf("%s: %v", msg, id)
			http.Error(w, msg, http.StatusNotFound)
			return
		}

		msg := "Failed to delete reservation"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Reservation deleted: %v", id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting reservation")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		msg := "Missing reservation ID"
		s.Logger.Error.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	reservation, err := s.MongoDB.GetReservation(id)
	if err != nil {
		msg := "Failed to get reservation from MongoDB"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(reservation)
	if err != nil {
		msg := "Failed to encode reservation to JSON"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User found: %v", reservation.ReservationID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (s *Server) getAllReservations(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all reservations")
	reservations, err := s.MongoDB.GetAll()
	if err != nil {
		msg := "Failed to get reservations"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(reservations)
	if err != nil {
		msg := "Failed to encode reservations to JSON"
		s.Logger.Error.Printf("%s: %v", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Reservations found: %v", len(reservations))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
