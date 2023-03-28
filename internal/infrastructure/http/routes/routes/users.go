package routes

import (
	"api/internal/infrastructure/http/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:          "/usuarios",
		Method:       http.MethodPost,
		Function:     controllers.CreateUser,
		RequiresAuth: false,
	},
	{
		URI:          "/usuarios",
		Method:       http.MethodGet,
		Function:     controllers.GetUsers,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}",
		Method:       http.MethodGet,
		Function:     controllers.GetUser,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}",
		Method:       http.MethodPut,
		Function:     controllers.UpdateUser,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteUser,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}/seguir",
		Method:       http.MethodPost,
		Function:     controllers.FollowUser,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}/parardeseguir",
		Method:       http.MethodDelete,
		Function:     controllers.StopFollowingUser,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}/seguidores",
		Method:       http.MethodGet,
		Function:     controllers.GetFollowers,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}/seguindo",
		Method:       http.MethodGet,
		Function:     controllers.GetFollowing,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}/atualizar-senha",
		Method:       http.MethodPost,
		Function:     controllers.UpdatePassword,
		RequiresAuth: true,
	},
}
