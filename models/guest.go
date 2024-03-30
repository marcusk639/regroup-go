package guestmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Guest struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}
