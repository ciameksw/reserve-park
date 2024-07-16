package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetServer(host string, port string) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/users", addUser).Methods("POST")
	router.HandleFunc("/users", editUser).Methods("PUT")
	router.HandleFunc("/users", deleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")

	return &http.Server{
		Addr:    host + ":" + port,
		Handler: router,
	}
}
