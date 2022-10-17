package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// TODO move variables to config file
const connStr = "mongodb://localhost:27017"

var Client *mongo.Client = ConnectDB()

// ConnectDB is a function to open connection with database
func ConnectDB() *mongo.Client {
	// Create new mongodb client
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatal(err)
	}

	// TODO think about its use
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongodb")

	return client
}

// OpenCollection is a function to make connection with database and open collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Get specified collection from database
	collection := client.Database("bananzaDB").Collection(collectionName)

	return collection
}
