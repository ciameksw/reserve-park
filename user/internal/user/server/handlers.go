package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	m "github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type addInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var input addInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := m.User{
		UserID:          uuid.NewString(),
		Username:        input.Username,
		Email:           input.Email,
		TotalMoneySpent: 0,
	}

	_, err = m.GetCollection().InsertOne(r.Context(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type editInput struct {
	UserID          string  `json:"user_id"`
	Username        string  `json:"username,omitempty"`
	Email           string  `json:"email,omitempty"`
	TotalMoneySpent float64 `json:"total_money_spent,omitempty"`
}

func editUser(w http.ResponseWriter, r *http.Request) {
	var input editInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.M{"user_id": bson.M{"$eq": input.UserID}}
	update := bson.M{}
	if input.Username != "" {
		update["username"] = input.Username
	}
	if input.Email != "" {
		update["email"] = input.Email
	}
	if input.TotalMoneySpent != 0 {
		update["total_money_spent"] = input.TotalMoneySpent
	}

	res, err := m.GetCollection().UpdateOne(r.Context(), filter, bson.M{"$set": update})
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	if res.MatchedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"user_id": bson.M{"$eq": id}}
	res := m.GetCollection().FindOneAndDelete(r.Context(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"user_id": bson.M{"$eq": id}}
	res := m.GetCollection().FindOne(r.Context(), filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	var user m.User
	err := res.Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	res, err := m.GetCollection().Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	var users []m.User
	err = res.All(r.Context(), &users)
	if err != nil {
		http.Error(w, "Failed to decode users", http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
