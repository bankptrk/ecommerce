package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var ProductCollection *mongo.Collection

func SetupDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	Client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect Database Successful")
	UserCollection = Client.Database("ecommerce").Collection("users")
	ProductCollection = Client.Database("ecommerce").Collection("products")
}

func GetUserCollection() *mongo.Collection {
	return UserCollection
}

func GetProductCollection() *mongo.Collection {
	return ProductCollection
}
