package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ReservationID string             `json:"reservation_id" bson:"reservation_id"`
	UserID        string             `json:"user_id" bson:"user_id"`
	SpotID        string             `json:"spot_id" bson:"spot_id"`
	Start         time.Time          `json:"start" bson:"start"`
	End           time.Time          `json:"end" bson:"end"`
	Canceled      bool               `json:"canceled" bson:"canceled"`
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
