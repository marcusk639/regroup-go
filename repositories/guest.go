package guestrepository

import (
	"context"
	"errors"
	models "regroup-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GuestRepository struct {
	DB         *mongo.Client
	Collection *mongo.Collection
}

func (g *GuestRepository) GetGuests() (*[]models.Guest, error) {
	var guests []models.Guest

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	cursor, err := g.Collection.Find(ctx, gin.H{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &guests)
	if err != nil {
		return nil, err
	}

	return &guests, nil
}

func (g *GuestRepository) GetGuest(id string) (*models.Guest, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result := g.Collection.FindOne(ctx, gin.H{"_id": objectID})

	var guest models.Guest

	err = result.Decode(&guest)
	if err != nil {
		return nil, err
	}

	return &guest, nil
}

func (g *GuestRepository) SaveGuest(guest models.Guest) (*models.Guest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	guest.ID = primitive.NewObjectID()

	_, err := g.Collection.InsertOne(ctx, guest)

	if err != nil {
		return nil, err
	}

	return &guest, nil
}

func (g *GuestRepository) UpdateGuest(id string, guest models.Guest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := g.Collection.ReplaceOne(ctx, gin.H{"_id": objectID}, guest)

	if result.MatchedCount == 0 {
		return errors.New("guest not found")
	}

	if err != nil {
		return nil
	}

	return nil
}

func (g *GuestRepository) DeleteGuest(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := g.Collection.DeleteOne(ctx, gin.H{"_id": objectID})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("guest not found")
	}

	return nil
}
