package db

import (
	"bishamon/config"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Client

// ! IMPORT AND INIT gotdotenv Load() before invoke this function !
func CreateMongoClient() (*mongo.Client, error) {
	conf := config.New()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conf.Mongo.CONNECT_URI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB", err) //log
	}
	fmt.Println("Successfully connected to MongoDB!") //log

	return client, nil
}

func InitMongo() error {
	client, err := CreateMongoClient()
	if err != nil {
		return err
	}

	Mongo = client
	return nil
}

func CloseMongo() error {
	return Mongo.Disconnect(context.Background())
}
