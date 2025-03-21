package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleType string

const (
	RoleAdmin RoleType = "admin"
	RoleUser  RoleType = "user"
)

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       string             `json:"user_id" bson:"user_id" validate:"required"`
	Username     string             `json:"username" bson:"username" validate:"required,min=3,max=30"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	PasswordHash string             `json:"password_hash" bson:"password_hash" validate:"required"`
	Role         RoleType           `json:"role" bson:"role" validate:"required,oneof=admin user"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at" validate:"required"`
}

func (m *MongoDB) AddUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.Collection.InsertOne(ctx, user)
	return err
}

func (m *MongoDB) GetUserByUsernameOrEmail(username, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}

	var user User
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *MongoDB) EditUser(input User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": bson.M{"$eq": input.UserID}}

	res := m.Collection.FindOneAndReplace(ctx, filter, input)
	return res.Err()
}

func (m *MongoDB) DeleteUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": bson.M{"$eq": userID}}

	res := m.Collection.FindOneAndDelete(ctx, filter)
	return res.Err()
}

func (m *MongoDB) GetUser(userID string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": bson.M{"$eq": userID}}

	var user User
	err := m.Collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (m *MongoDB) GetAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	err = cursor.All(ctx, &users)
	return users, err
}
