package guestservice

import (
	"context"
	"net/http"
	models "regroup-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GuestService struct {
	DB *mongo.Client
}

func GetCollection(db *mongo.Client) *mongo.Collection {
	return db.Database("regroup").Collection("guests")
}

func (service *GuestService) SaveGuest(c *gin.Context) {
	var guest models.Guest

	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := GetCollection(service.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, err := collection.InsertOne(ctx, guest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error saving guest"})
	}

	c.JSON(http.StatusOK, guest)
}

func (service *GuestService) GetGuests(c *gin.Context) {
	var guests []models.Guest

	collection := GetCollection(service.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, gin.H{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting guests"})
	}

	defer cursor.Close(ctx)

	unmarshalError := cursor.All(ctx, &guests)
	if unmarshalError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting guests"})
	}

	c.JSON(http.StatusOK, guests)
}

func (service *GuestService) GetGuest(c *gin.Context) {
	id := c.Param("id")

	objectID, objectIdErr := primitive.ObjectIDFromHex(id)
	if objectIdErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	collection := GetCollection(service.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var guest models.Guest
	result := collection.FindOne(ctx, gin.H{"_id": objectID})
	err := result.Decode(&guest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting guest"})
	}

	c.JSON(http.StatusOK, guest)
}

func (service *GuestService) DeleteGuest(c *gin.Context) {
	id := c.Param("id")

	objectID, objectIdErr := primitive.ObjectIDFromHex(id)
	if objectIdErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	collection := GetCollection(service.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := collection.DeleteOne(ctx, gin.H{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting guest"})
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guest not found"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guest deleted"})
}

func (service *GuestService) UpdateGuest(c *gin.Context) {
	id := c.Param("id")

	objectID, objectIdErr := primitive.ObjectIDFromHex(id)
	if objectIdErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	var guest models.Guest

	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := GetCollection(service.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := collection.ReplaceOne(ctx, gin.H{"_id": objectID}, guest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating guest"})
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guest not found"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guest updated"})
}
