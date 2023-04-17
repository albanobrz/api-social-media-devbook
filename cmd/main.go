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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}