package server

import (
	"log"
	"net/http"

	"github.com/ciameksw/reserve-park/user/internal/user/config"
	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/gorilla/mux"
)

type Server struct {
	Config  *config.Config
	MongoDB *mongodb.MongoDB
}

func NewServer(cfg *config.Config, db *mongodb.MongoDB) *Server {
	return &Server{
		Config:  cfg,
		MongoDB: db,
	}
}

func (s *Server) Start() {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.addUser).Methods("POST")
	r.HandleFunc("/users", s.editUser).Methods("PUT")
	r.HandleFunc("/users/{id}", s.deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", s.getUser).Methods("GET")
	r.HandleFunc("/users", s.getAllUsers).Methods("GET")

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	log.Printf("Server started at %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
