package mongodb

import (
	"context"
	"time"

	"github.com/benweissmann/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMock() (*MongoDB, error) {
	mongoServer, err := memongo.Start("4.2.8")
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(mongoServer.URI())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database("mock").Collection("mock")

	return &MongoDB{Collection: collection}, nil
}
