package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	godotenv.Load("../.env")

	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	dbUri := os.Getenv("DB_URI")
	if dbUri == "" {
		log.Fatal("DB_URI environment variable is required")
	}

	fmt.Println("Connecting to MongoDB...")
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("Likky")
	collection := db.Collection("Users")

	// Find the user and show raw document
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Raw document from MongoDB:")
		for key, value := range raw {
			fmt.Printf("  %s: %T\n", key, value)
		}
		fmt.Printf("Full document: %+v\n", raw)
		break // Just show first document
	}
}
