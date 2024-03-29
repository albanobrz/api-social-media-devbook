package main

import (
	"api/internal/infrastructure/config"
	"api/internal/infrastructure/database"
	"api/internal/infrastructure/http/middlewares"
	router "api/internal/infrastructure/http/router/mux/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func configurateRoutes(r *mux.Router, db *mongo.Database) *mux.Router {
	routes := []router.Route{}

	usersRoutes := router.ConfigUsersRoutes(db)
	postsRoutes := router.ConfigPostsRoutes(db)
	loginRoute := router.ConfigLoginRoutes(db)

	routes = append(routes, usersRoutes...)
	routes = append(routes, postsRoutes...)
	routes = append(routes, loginRoute)

	for _, route := range routes {
		if route.RequiresAuth {
			r.HandleFunc(route.URI,
				middlewares.Logger(
					middlewares.Authenticate(route.Controller),
				),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, route.Controller).Methods(route.Method)
		}
	}

	return r
}

func main() {
	config.Load()

	mongo, err := database.Connect()
	if err != nil {
		panic(fmt.Errorf("Could not connect on mongoDB: %s", err))
	}

	r := mux.NewRouter()

	configRoutes := configurateRoutes(r, mongo)

	var PORT = fmt.Sprintf(":%v", config.Port)

	fmt.Printf("Listening on PORT %v...\n", config.Port)
	log.Fatal(http.ListenAndServe(PORT, configRoutes))
}
