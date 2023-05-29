package routes

import (
	"api/internal/infrastructure/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI          string
	Method       string
	Function     func(http.ResponseWriter, *http.Request)
	RequiresAuth bool
}

func Config(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, LoginRoutes)
	routes = append(routes, PostsRoutes...)

	for _, route := range routes {
		if route.RequiresAuth {
			r.HandleFunc(route.URI, middlewares.Logger(middlewares.Authenticate(route.Function))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Method)
		}
	}
	return r
}
