package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Spot struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SpotID       string             `json:"spot_id" bson:"spot_id"`
	Latitude     float64            `json:"latitude" bson:"latitude"`
	Longitude    float64            `json:"longitude" bson:"longitude"`
	PricePerHour float64            `json:"price_per_hour" bson:"price_per_hour"`
}

func (m *MongoDB) AddSpot(spot Spot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.Collection.InsertOne(ctx, spot)
	return err
}

func (m *MongoDB) EditSpot(input Spot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"spot_id": bson.M{"$eq": input.SpotID}}

	res := m.Collection.FindOneAndReplace(ctx, filter, input)
	return res.Err()
}

func (m *MongoDB) DeleteSpot(spotID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"spot_id": bson.M{"$eq": spotID}}

	res := m.Collection.FindOneAndDelete(ctx, filter)
	return res.Err()
}

func (m *MongoDB) GetSpot(spotID string) (Spot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"spot_id": bson.M{"$eq": spotID}}

	var spot Spot
	err := m.Collection.FindOne(ctx, filter).Decode(&spot)
	return spot, err
}

func (m *MongoDB) GetAll() ([]Spot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var spots []Spot
	err = cursor.All(ctx, &spots)
	return spots, err
}
