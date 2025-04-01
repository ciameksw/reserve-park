package server

import (
	"net/http"

	"github.com/ciameksw/reserve-park/facade/internal/facade/config"
	"github.com/ciameksw/reserve-park/facade/internal/facade/logger"
	"github.com/ciameksw/reserve-park/facade/internal/facade/services/reservation"
	"github.com/ciameksw/reserve-park/facade/internal/facade/services/spot"
	"github.com/ciameksw/reserve-park/facade/internal/facade/services/user"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Server struct {
	Logger             *logger.Logger
	Config             *config.Config
	UserService        *user.UserService
	SpotService        *spot.SpotService
	ReservationService *reservation.ReservationService
	Validator          *validator.Validate
}

func NewServer(log *logger.Logger,
	cfg *config.Config,
	usr *user.UserService,
	spt *spot.SpotService,
	rsrv *reservation.ReservationService) *Server {
	return &Server{
		Logger:             log,
		Config:             cfg,
		UserService:        usr,
		SpotService:        spt,
		ReservationService: rsrv,
		Validator:          validator.New(),
	}
}

func (s *Server) Start() {
	r := mux.NewRouter()

	r.HandleFunc("/login", s.login).Methods("POST")

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	s.Logger.Info.Printf("Server started at %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		s.Logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
