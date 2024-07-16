package server

import "net/http"

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Add user"))
}

func editUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Edit user"))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get user"))
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all users"))
}
