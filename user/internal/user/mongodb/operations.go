package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	Username        string             `json:"username" bson:"username"`
	Email           string             `json:"email" bson:"email"`
	TotalMoneySpent float64            `json:"total_money_spent" bson:"total_money_spent"`
}

func (m *MongoDB) AddUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.Collection.InsertOne(ctx, user)
	return err
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
