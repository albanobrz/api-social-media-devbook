package routes

import (
	"api/internal/infrastructure/http/controllers"
	"net/http"
)

var LoginRoutes = Route{
	URI:          "/login",
	Method:       http.MethodPost,
	Function:     controllers.Login,
	RequiresAuth: false,
}
