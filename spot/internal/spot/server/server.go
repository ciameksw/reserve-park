package server

import (
	"net/http"

	"github.com/ciameksw/reserve-park/spot/internal/spot/config"
	"github.com/ciameksw/reserve-park/spot/internal/spot/logger"
	"github.com/ciameksw/reserve-park/spot/internal/spot/mongodb"
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

	r.HandleFunc("/spots/price", s.getPrice).Methods("GET")

	r.HandleFunc("/spots", s.addSpot).Methods("POST")
	r.HandleFunc("/spots", s.editSpot).Methods("PUT")
	r.HandleFunc("/spots/{id}", s.deleteSpot).Methods("DELETE")
	r.HandleFunc("/spots/{id}", s.getSpot).Methods("GET")
	r.HandleFunc("/spots", s.getAllSpots).Methods("GET")

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	s.Logger.Info.Printf("Server started at %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		s.Logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
