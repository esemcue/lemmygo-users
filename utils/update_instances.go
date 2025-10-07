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

type credentials map[string]string

type Instance struct {
	URL         string      `bson:"url" json:"url"`
	Credentials credentials `bson:"credentials" json:"credentials"`
}

type User struct {
	Email     string              `bson:"_id" json:"Email"`
	Password  string              `bson:"password" json:"Password"`
	Instances map[string]Instance `bson:"instances" json:"Instances"`
}

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

	// Define sample instances
	sampleInstances := map[string]Instance{
		"lemmy.world": {
			URL: "https://lemmy.world",
			Credentials: credentials{
				"username": "demo_user",
				"token":    "demo_token_123",
			},
		},
		"beehaw.org": {
			URL: "https://beehaw.org",
			Credentials: credentials{
				"username": "test_user",
				"password": "secure_password",
			},
		},
		"lemmy.ml": {
			URL: "https://lemmy.ml",
			Credentials: credentials{
				"api_key": "api_key_456",
			},
		},
		"sh.itjust.works": {
			URL: "https://sh.itjust.works",
			Credentials: credentials{
				"username": "worker",
				"jwt":      "jwt_token_789",
			},
		},
	}

	// First, let's see if there are any existing users
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d existing users\n", len(users))

	if len(users) == 0 {
		fmt.Println("No existing users found. Creating a demo user with instances...")

		// Create a demo user with instances
		demoUser := User{
			Email:     "demo@example.com",
			Password:  "$2a$10$demo.hash.here", // This would be a real bcrypt hash
			Instances: sampleInstances,
		}

		result, err := collection.InsertOne(context.TODO(), demoUser)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Created demo user with ID: %v\n", result.InsertedID)
	} else {
		// Update existing users to add instances
		for _, user := range users {
			fmt.Printf("Updating user: %s\n", user.Email)

			// Merge existing instances with sample instances
			updatedInstances := make(map[string]Instance)

			// Keep existing instances
			if user.Instances != nil {
				for k, v := range user.Instances {
					updatedInstances[k] = v
				}
			}

			// Add sample instances (don't overwrite existing ones)
			for k, v := range sampleInstances {
				if _, exists := updatedInstances[k]; !exists {
					updatedInstances[k] = v
				}
			}

			// Update the user document
			filter := bson.D{{Key: "_id", Value: user.Email}}
			update := bson.D{{Key: "$set", Value: bson.D{{Key: "instances", Value: updatedInstances}}}}

			result, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Printf("Error updating user %s: %v\n", user.Email, err)
				continue
			}

			fmt.Printf("Updated user %s: matched %d, modified %d\n",
				user.Email, result.MatchedCount, result.ModifiedCount)
		}
	}

	fmt.Println("Database update completed successfully!")

	// Show final state
	fmt.Println("\nFinal user data:")
	cursor, err = collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var finalUsers []User
	if err = cursor.All(context.TODO(), &finalUsers); err != nil {
		log.Fatal(err)
	}

	for _, user := range finalUsers {
		fmt.Printf("\nUser: %s\n", user.Email)
		fmt.Printf("Instances (%d):\n", len(user.Instances))
		for name, instance := range user.Instances {
			fmt.Printf("  - %s: %s\n", name, instance.URL)
			fmt.Printf("    Credentials: %d items\n", len(instance.Credentials))
			for k, v := range instance.Credentials {
				fmt.Printf("      %s: %s\n", k, v)
			}
		}
	}
}
