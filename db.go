package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DbURI := fmt.Sprintf("mongodb://%s:%s", DBHost, DBPort)
	clientOptions := options.Client().ApplyURI(DbURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}
