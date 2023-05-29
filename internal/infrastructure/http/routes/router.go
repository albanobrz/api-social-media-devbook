package router

import (
	"api/internal/infrastructure/http/routes/routes"

	"github.com/gorilla/mux"
)

func Create() *mux.Router {
	r := mux.NewRouter()
	return routes.Config(r)
}
