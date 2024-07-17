package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	Username        string             `json:"username" bson:"username"`
	Email           string             `json:"email" bson:"email"`
	TotalMoneySpent float64            `json:"total_money_spent" bson:"total_money_spent"`
}

type MongoDB struct {
	Collection *mongo.Collection
}

func Connect(uri string, db string) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDB{Collection: client.Database(db).Collection(db)}, nil
}

func (m *MongoDB) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.Collection.Database().Client().Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
