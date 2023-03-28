package router

import (
	"api/internal/infrastructure/http/routes/routes"

	"github.com/gorilla/mux"
)

// Gerar vai retornar um router com as rotas configuradas
func Create() *mux.Router {
	r := mux.NewRouter()
	return routes.Config(r)
}
