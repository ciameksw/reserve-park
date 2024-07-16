package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoCollection *mongo.Collection

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

	collection := client.Database(db).Collection(db)

	mongoCollection = collection
	fmt.Printf("Connected to MongoDB, db: %s collection: %s\n", mongoCollection.Database().Name(), mongoCollection.Name())
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
