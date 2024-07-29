package database

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

func ConnectDB() {
	MONGODB_URL := os.Getenv("MONGODB_URL")

	if MONGODB_URL == "" {
		log.Fatal("Error getting MONGODB_URL from the environment variable")
	}

	clientOptions := options.Client().ApplyURI(MONGODB_URL)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	DB = client.Database("golang_ecomerce")

	fmt.Println("Connected to MongoDB Atlas")
}

var JWTSecret = []byte(string(os.Getenv("SECRET_KEY")))
