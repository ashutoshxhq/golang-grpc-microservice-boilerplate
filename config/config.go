package config

import (
	"context"
	"log"

	"egnite.app/microservices/user/config/development"
	"egnite.app/microservices/user/config/production"
	"egnite.app/microservices/user/config/staging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Database is global database variable
	Database *mongo.Client
	// Environment is global varible for setting the env
	Environment string = "development"
	// URI is mongo connection string
	URI string
)

// InitialiseEnvironment is a database initialiser
func InitialiseEnvironment() {

	if Environment == "development" {
		URI = development.URI
	} else if Environment == "production" {
		URI = production.URI
	} else if Environment == "staging" {
		URI = staging.URI
	} else {
		URI = development.URI
	}

	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	Database = client
}
