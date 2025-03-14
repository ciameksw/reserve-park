package server

import (
	"net/http"

	"github.com/ciameksw/reserve-park/reservation/internal/reservation/config"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/logger"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/mongodb"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Server struct {
	Logger    *logger.Logger
	Config    *config.Config
	MongoDB   *mongodb.MongoDB
	Validator *validator.Validate
}

func NewServer(log *logger.Logger, cfg *config.Config, db *mongodb.MongoDB) *Server {
	return &Server{
		Logger:    log,
		Config:    cfg,
		MongoDB:   db,
		Validator: validator.New(),
	}
}

func (s *Server) Start() {
	r := mux.NewRouter()

	r.HandleFunc("/reservations", s.addReservation).Methods("POST")
	r.HandleFunc("/reservations", s.editReservation).Methods("PUT")
	r.HandleFunc("/reservations/{id}", s.deleteReservation).Methods("DELETE")
	r.HandleFunc("/reservations/{id}", s.getReservation).Methods("GET")
	r.HandleFunc("/reservations", s.getAllReservations).Methods("GET")

	r.HandleFunc("/reservations/user/{id}", s.getUserReservations).Methods("GET")
	r.HandleFunc("/reservations/spot/{id}", s.getSpotReservations).Methods("GET")

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	s.Logger.Info.Printf("Server started at %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		s.Logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
