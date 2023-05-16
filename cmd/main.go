package main

import (
	"api/internal/infrastructure/config"
	router "api/internal/infrastructure/http/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Load()
	r := router.Create()

	fmt.Printf("Listening in port %d", config.Port)

	// clientOptions := options.Client().ApplyURI(config.ConnectionDBStringMongo)
	// client, err := mongo.Connect(context.Background(), clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = client.Ping(context.Background(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Connected to mongoDB")

	// db, err := database.ConnectMongo()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
