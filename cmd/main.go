package main

import (
	"api/internal/infrastructure/config"
	router "api/internal/infrastructure/http/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.Load()
	r := router.Create()

	fmt.Printf("Listening in port %d", config.Port)

	ConnectionDBString := fmt.Sprintf("mongodb://%s:%s@172.19.0.2:27017/",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	clientOptions := options.Client().ApplyURI(ConnectionDBString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
