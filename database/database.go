package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() *mongo.Client {
	connectionString := "mongodb+srv://marcusk639:MpXnDflRTFLtEdEp@regroup.b0otx1w.mongodb.net/?retryWrites=true&w=majority&appName=Regroup"
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	return client
}

func GetCollection(db *mongo.Client, name string) *mongo.Collection {
	return db.Database("regroup").Collection(name)
}
