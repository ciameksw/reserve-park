package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoCollection *mongo.Collection

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	Username        string             `json:"username" bson:"username"`
	Email           string             `json:"email" bson:"email"`
	TotalMoneySpent float64            `json:"total_money_spent" bson:"total_money_spent"`
}

func Connect(uri string, db string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	c := client.Database(db).Collection(db)
	mongoCollection = c
	fmt.Printf("Connected to MongoDB, db: %s collection: %s\n", c.Database().Name(), c.Name())
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := mongoCollection.Database().Client().Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func GetCollection() *mongo.Collection {
	return mongoCollection
}
