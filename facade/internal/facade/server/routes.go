package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) addUserRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("/register", s.register).Methods("POST")
	userRouter.HandleFunc("/login", s.login).Methods("POST")

	// Admin routes
	userRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.getAllUsers))).Methods("GET")
	userRouter.Handle("/role", s.authorize(RoleAdmin, http.HandlerFunc(s.editUsersRole))).Methods("PATCH")

	// User routes
	userRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.editUser))).Methods("PATCH")
	userRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getUserByID))).Methods("GET")
	userRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.deleteUserByID))).Methods("DELETE")
}

func (s *Server) addSpotRoutes(r *mux.Router) {
	spotRouter := r.PathPrefix("/spots").Subrouter()

	// User routes
	spotRouter.Handle("/available", s.authorize(RoleUser, http.HandlerFunc(s.register))).Methods("GET") // TODO: Implement real handler
	spotRouter.Handle("/price", s.authorize(RoleUser, http.HandlerFunc(s.getSpotPrice))).Methods("GET")
	spotRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.getAllSpots))).Methods("GET")
	spotRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getSpotByID))).Methods("GET")

	// Admin routes
	spotRouter.Handle("/{id}", s.authorize(RoleAdmin, http.HandlerFunc(s.deleteSpotByID))).Methods("DELETE")
	spotRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.addSpot))).Methods("POST")
	spotRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.editSpot))).Methods("PATCH")
}

func (s *Server) addReservationRoutes(r *mux.Router) {
	reservationRouter := r.PathPrefix("/reservations").Subrouter()

	// Admin routes
	reservationRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.getAllReservations))).Methods("GET")
	reservationRouter.Handle("/spot/{id}", s.authorize(RoleAdmin, http.HandlerFunc(s.getReservationsBySpot))).Methods("GET")
	reservationRouter.Handle("/{id}", s.authorize(RoleAdmin, http.HandlerFunc(s.deleteReservationByID))).Methods("DELETE")

	// User routes
	reservationRouter.Handle("/user/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getReservationsByUser))).Methods("GET")
	reservationRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getReservationByID))).Methods("GET")
	reservationRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.addReservation))).Methods("POST")
	reservationRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.editReservation))).Methods("PATCH")
	reservationRouter.Handle("/cancel/{id}", s.authorize(RoleUser, http.HandlerFunc(s.cancelReservation))).Methods("PATCH")
}
