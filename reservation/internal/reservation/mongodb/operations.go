package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusType string

const (
	StatusValid    StatusType = "valid"
	StatusCanceled StatusType = "canceled"
)

type Reservation struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ReservationID string             `json:"reservation_id" bson:"reservation_id" validate:"required"`
	UserID        string             `json:"user_id" bson:"user_id" validate:"required"`
	SpotID        string             `json:"spot_id" bson:"spot_id" validate:"required"`
	StartTime     time.Time          `json:"start_time" bson:"start_time" validate:"required"`
	EndTime       time.Time          `json:"end_time" bson:"end_time" validate:"required"`
	Status        StatusType         `json:"status" bson:"status" validate:"required,oneof=valid canceled"`
	PricePaid     float64            `json:"price_paid" bson:"price_paid" validate:"required,gt=0"`
	CreateAt      time.Time          `json:"created_at" bson:"created_at" validate:"required"`
}

func (m *MongoDB) AddReservation(reservation Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.Collection.InsertOne(ctx, reservation)
	return err
}

func (m *MongoDB) EditReservation(input Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"reservation_id": bson.M{"$eq": input.ReservationID}}

	res := m.Collection.FindOneAndReplace(ctx, filter, input)
	return res.Err()
}

func (m *MongoDB) DeleteReservation(reservationID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"reservation_id": bson.M{"$eq": reservationID}}

	res := m.Collection.FindOneAndDelete(ctx, filter)
	return res.Err()
}

func (m *MongoDB) GetReservation(reservationID string) (Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"reservation_id": bson.M{"$eq": reservationID}}

	var reservation Reservation
	err := m.Collection.FindOne(ctx, filter).Decode(&reservation)
	return reservation, err
}

func (m *MongoDB) GetAll() ([]Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reservations []Reservation
	err = cursor.All(ctx, &reservations)
	return reservations, err
}

type parameterInput string

const (
	ByUserID parameterInput = "user_id"
	BySpotID parameterInput = "spot_id"
)

func (m *MongoDB) GetReservationsBy(parameter parameterInput, id string) ([]Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{string(parameter): bson.M{"$eq": id}}

	cursor, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var reservations []Reservation
	err = cursor.All(ctx, &reservations)
	return reservations, err
}

type AvailabilityInput struct {
	SpotIDs   []string  `json:"spot_ids" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

func (m *MongoDB) CheckAvailability(input AvailabilityInput) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"spot_id":    bson.M{"$in": input.SpotIDs},
		"start_time": bson.M{"$lt": input.EndTime},
		"end_time":   bson.M{"$gt": input.StartTime},
	}

	cursor, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var reservations []Reservation
	err = cursor.All(ctx, &reservations)
	if err != nil {
		return nil, err
	}

	// Create map of unavailable spots
	unavailableSpots := make(map[string]struct{})
	for _, res := range reservations {
		unavailableSpots[res.SpotID] = struct{}{}
	}

	var availableSpots []string
	for _, spotID := range input.SpotIDs {
		_, found := unavailableSpots[spotID]
		if !found {
			availableSpots = append(availableSpots, spotID)
		}
	}

	return availableSpots, nil
}
