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
	spotRouter.Handle("/available", s.authorize(RoleUser, http.HandlerFunc(s.register))).Methods("GET") // TODO
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
	reservationRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.getAllReservations))).Methods("GET")              // GIT get all reservations
	reservationRouter.Handle("/spot/{id}", s.authorize(RoleAdmin, http.HandlerFunc(s.getReservationsBySpot))).Methods("GET") // GIT get reservations for a spot (admin)
	reservationRouter.Handle("/{id}", s.authorize(RoleAdmin, http.HandlerFunc(s.deleteReservationByID))).Methods("DELETE")   // GIT delete reservation (admin)

	// User routes
	reservationRouter.Handle("/user/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getReservationsByUser))).Methods("GET") // GIT get reservations for user if it is him (admin for everybody)
	reservationRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getReservationByID))).Methods("GET")         // GIT get reservation if its users (admin for everybody)
	reservationRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.addReservation))).Methods("POST")                 // GIT create reservation if for user (admin for all)
	reservationRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.editReservation))).Methods("PATCH")               // GIT edit reservation if his (admin all)
	reservationRouter.Handle("/cancel/{id}", s.authorize(RoleUser, http.HandlerFunc(s.cancelReservation))).Methods("PATCH") // cancel if its users (admin for all)
}
