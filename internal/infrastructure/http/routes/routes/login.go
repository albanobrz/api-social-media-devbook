package routes

import (
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/controllers"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigLoginRoutes(db *mongo.Database) Route {

	repository := repositories.NewUsersRepository(db)

	controllers := controllers.NewUsersController(*repository)

	var LoginRoutes = Route{
		URI:          "/login",
		Method:       http.MethodPost,
		Function:     controllers.Login,
		RequiresAuth: false,
	}
	return LoginRoutes
}
