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
	UserID    string       `json:"user_id"`
	SpotID    string       `json:"spot_id"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	Status    m.StatusType `json:"status"`
	PricePaid float64      `json:"price_paid"`
	CreateAt  time.Time    `json:"created_at"`
}

func (s *Server) addReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding reservation")
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	data := m.Reservation{
		ReservationID: uuid.NewString(),
		UserID:        input.UserID,
		SpotID:        input.SpotID,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		Status:        input.Status,
		PricePaid:     input.PricePaid,
		CreateAt:      input.CreateAt,
	}

	// TODO: should we check if the spot is available here?

	err = s.MongoDB.AddReservation(data)
	if err != nil {
		s.handleError(w, "Failed to add reservation to MongoDB", err, http.StatusInternalServerError)
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
		s.handleError(w, "Failed to decode request body", err, http.StatusBadRequest)
		return
	}

	err = s.MongoDB.EditReservation(input)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Reservation not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to edit reservation in MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Reservation edited: %v", input.ReservationID)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteReservation(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting reservation")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing reservation ID", nil, http.StatusBadRequest)
		return
	}

	err := s.MongoDB.DeleteReservation(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Reservation not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to delete reservation", err, http.StatusInternalServerError)
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
		s.handleError(w, "Missing reservation ID", nil, http.StatusBadRequest)
		return
	}

	reservation, err := s.MongoDB.GetReservation(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Reservation not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to get reservation from MongoDB", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("Reservation found: %v", reservation.ReservationID)
	s.writeJSON(w, reservation, http.StatusOK)
}

func (s *Server) getAllReservations(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting all reservations")
	reservations, err := s.MongoDB.GetAll()
	if err != nil {
		s.handleError(w, "Failed to get reservations", err, http.StatusInternalServerError)
		return
	}

	if len(reservations) == 0 {
		s.Logger.Info.Println("No reservations found")
		s.writeJSON(w, []m.Reservation{}, http.StatusOK)
		return
	}

	s.Logger.Info.Printf("Reservations found: %v documents", len(reservations))
	s.writeJSON(w, reservations, http.StatusOK)
}

func (s *Server) getUserReservations(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting reservations by userID")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing user ID", nil, http.StatusBadRequest)
		return
	}

	reservations, err := s.MongoDB.GetReservationsBy("user_id", id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Reservations not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to get reservations from MongoDB", err, http.StatusInternalServerError)
		return
	}

	if len(reservations) == 0 {
		s.Logger.Info.Println("No reservations found")
		s.writeJSON(w, []m.Reservation{}, http.StatusOK)
		return
	}

	s.Logger.Info.Printf("Reservations found: %v documents", len(reservations))
	s.writeJSON(w, reservations, http.StatusOK)
}

func (s *Server) getSpotReservations(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting reservations by spotID")
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		s.handleError(w, "Missing spot ID", nil, http.StatusBadRequest)
		return
	}

	reservations, err := s.MongoDB.GetReservationsBy("spot_id", id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.handleError(w, "Reservations not found", err, http.StatusNotFound)
			return
		}

		s.handleError(w, "Failed to get reservations from MongoDB", err, http.StatusInternalServerError)
		return
	}

	if len(reservations) == 0 {
		s.Logger.Info.Println("No reservations found")
		s.writeJSON(w, []m.Reservation{}, http.StatusOK)
		return
	}

	s.Logger.Info.Printf("Reservations found: %v documents", len(reservations))
	s.writeJSON(w, reservations, http.StatusOK)
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
