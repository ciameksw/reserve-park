package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	Username        string             `json:"username" bson:"username"`
	Email           string             `json:"email" bson:"email"`
	TotalMoneySpent float64            `json:"total_money_spent" bson:"total_money_spent"`
}

func (m *MongoDB) AddUser(user User) error {
	_, err := m.Collection.InsertOne(context.Background(), user)
	return err
}

func (m *MongoDB) EditUser(input User) error {
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

	res, err := m.Collection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (m *MongoDB) DeleteUser(userID string) error {
	filter := bson.M{"user_id": bson.M{"$eq": userID}}
	res := m.Collection.FindOneAndDelete(context.Background(), filter)
	return res.Err()
}

func (m *MongoDB) GetUser(userID string) (User, error) {
	filter := bson.M{"user_id": bson.M{"$eq": userID}}
	var user User
	err := m.Collection.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}

func (m *MongoDB) GetAll() ([]User, error) {
	cursor, err := m.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []User
	err = cursor.All(context.Background(), &users)
	return users, err
}
