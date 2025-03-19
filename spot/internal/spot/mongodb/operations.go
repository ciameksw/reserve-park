package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SizeType string

const (
	SizeSmall  SizeType = "small"
	SizeMedium SizeType = "medium"
	SizeLarge  SizeType = "large"
)

type SpotType string

const (
	SpotTypeIndoor  SpotType = "indoor"
	SpotTypeOutdoor SpotType = "outdoor"
	SpotTypeEV      SpotType = "ev"
)

type Spot struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SpotID       string             `json:"spot_id" bson:"spot_id" validate:"required"`
	Latitude     float64            `json:"latitude" bson:"latitude" validate:"required"`
	Longitude    float64            `json:"longitude" bson:"longitude" validate:"required"`
	PricePerHour float64            `json:"price_per_hour" bson:"price_per_hour" validate:"required,gt=0"`
	Size         SizeType           `json:"size" bson:"size" validate:"required,oneof=small medium large"`
	Type         SpotType           `json:"type" bson:"type" validate:"required,oneof=indoor outdoor ev"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at" validate:"required"`
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

type GetPriceInput struct {
	SpotID    string    `json:"spot_id" bson:"spot_id" validate:"required"`
	StartTime time.Time `json:"start_time" bson:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" bson:"end_time" validate:"required"`
}

func (m *MongoDB) GetPrice(input GetPriceInput) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"spot_id": bson.M{"$eq": input.SpotID}}

	var spot Spot
	err := m.Collection.FindOne(ctx, filter).Decode(&spot)
	if err != nil {
		return 0, err
	}

	diff := input.EndTime.Sub(input.StartTime).Hours()
	price := diff * spot.PricePerHour
	price = float64(int(price*100)) / 100

	return price, nil
}
