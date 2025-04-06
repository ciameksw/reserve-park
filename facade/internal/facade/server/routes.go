package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) addUserRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("/register", s.register).Methods("POST")
	userRouter.HandleFunc("/login", s.login).Methods("POST")

	userRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.getAllUsers))).Methods("GET")
	userRouter.Handle("/role", s.authorize(RoleAdmin, http.HandlerFunc(s.editUsersRole))).Methods("PATCH")

	userRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.editUser))).Methods("PATCH")
	userRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getUserByID))).Methods("GET")
	userRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.deleteUserByID))).Methods("DELETE")
}

func (s *Server) addSpotRoutes(r *mux.Router) {
	spotRouter := r.PathPrefix("/spots").Subrouter()

	spotRouter.Handle("/price", s.authorize(RoleUser, http.HandlerFunc(s.getSpotPrice))).Methods("GET")
	spotRouter.Handle("", s.authorize(RoleUser, http.HandlerFunc(s.getAllSpots))).Methods("GET")
	spotRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getSpotByID))).Methods("GET")

	spotRouter.Handle("/{id}", s.authorize(RoleAdmin, http.HandlerFunc(s.deleteSpotByID))).Methods("DELETE")
	spotRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.addSpot))).Methods("POST")
	spotRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.editSpot))).Methods("PATCH")
}
