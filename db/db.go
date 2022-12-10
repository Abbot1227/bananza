package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// TODO move variables to config file
var username, password = "main", "admin"
var connStr = "mongodb+srv://" + username + ":" + password + "@cluster0.9rrlh4n.mongodb.net/?retryWrites=true&w=majority"

var Client = ConnectDB()

// getDBConfig is a function used to get username and password for mongdodb connection
func getDBConfig() (string, string) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGODB_NICKNAME"), os.Getenv("MONGODB_PASSWORD")
}

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
	collection := client.Database("bananza").Collection(collectionName)

	return collection
}
