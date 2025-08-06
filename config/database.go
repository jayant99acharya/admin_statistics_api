package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var Client *mongo.Client

func ConnectDatabase() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "admin_statistics"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	Client = client
	DB = client.Database(dbName)
	
	fmt.Println("Connected to MongoDB successfully!")
	
	// Create indexes for better performance
	createIndexes()
}

func createIndexes() {
	ctx := context.Background()
	collection := DB.Collection("transactions")

	// Index on createdAt for time-based queries
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int{"createdAt": 1},
	})
	if err != nil {
		log.Printf("Failed to create createdAt index: %v", err)
	}

	// Index on userId for user-specific queries
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int{"userId": 1},
	})
	if err != nil {
		log.Printf("Failed to create userId index: %v", err)
	}

	// Compound index on userId and createdAt
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int{"userId": 1, "createdAt": 1},
	})
	if err != nil {
		log.Printf("Failed to create compound index: %v", err)
	}

	// Index on roundId for round-based queries
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int{"roundId": 1},
	})
	if err != nil {
		log.Printf("Failed to create roundId index: %v", err)
	}

	// Index on type for filtering wagers/payouts
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int{"type": 1},
	})
	if err != nil {
		log.Printf("Failed to create type index: %v", err)
	}

	fmt.Println("Database indexes created successfully!")
}

func DisconnectDatabase() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		Client.Disconnect(ctx)
	}
}