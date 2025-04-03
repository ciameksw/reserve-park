package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) addUserRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/users").Subrouter()

	userRouter.HandleFunc("/login", s.login).Methods("POST")
	userRouter.Handle("", s.authorize(RoleAdmin, http.HandlerFunc(s.getAllUsers))).Methods("GET")
	userRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.getUserByID))).Methods("GET")
	userRouter.Handle("/{id}", s.authorize(RoleUser, http.HandlerFunc(s.deleteUserByID))).Methods("DELETE")
}
