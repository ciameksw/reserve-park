package server

import (
	"net/http"

	"github.com/ciameksw/reserve-park/user/internal/user/config"
	"github.com/ciameksw/reserve-park/user/internal/user/logger"
	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
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

	r.HandleFunc("/users/authorize", s.authorize).Methods("GET")

	r.HandleFunc("/users", s.addUser).Methods("POST")
	r.HandleFunc("/users", s.editUser).Methods("PATCH")
	r.HandleFunc("/users/{id}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", s.getUser).Methods("GET")
	r.HandleFunc("/users", s.getAllUsers).Methods("GET")

	r.HandleFunc("/users/login", s.login).Methods("POST")

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	s.Logger.Info.Printf("Server started at %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		s.Logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
