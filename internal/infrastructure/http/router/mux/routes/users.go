package routes

import (
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/controllers"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigUsersRoutes(db *mongo.Database) []Route {

	repository := repositories.NewUsersRepository(db)

	controllers := controllers.NewUsersController(repository)

	var userRoutes = []Route{
		{
			URI:          "/users",
			Method:       http.MethodPost,
			Controller:   controllers.CreateUser,
			RequiresAuth: false,
		},
		{
			URI:          "/users/",
			Method:       http.MethodGet,
			Controller:   controllers.GetAllUsers,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}",
			Method:       http.MethodGet,
			Controller:   controllers.GetUser,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}",
			Method:       http.MethodPut,
			Controller:   controllers.UpdateUser,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}",
			Method:       http.MethodDelete,
			Controller:   controllers.DeleteUser,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}/follow",
			Method:       http.MethodPost,
			Controller:   controllers.FollowUser,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}/unfollow",
			Method:       http.MethodDelete,
			Controller:   controllers.UnfollowUser,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}/followers",
			Method:       http.MethodGet,
			Controller:   controllers.GetFollowers,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}/following",
			Method:       http.MethodGet,
			Controller:   controllers.GetFollowing,
			RequiresAuth: true,
		},
		{
			URI:          "/users/{userID}/update-password",
			Method:       http.MethodPost,
			Controller:   controllers.UpdatePassword,
			RequiresAuth: true,
		},
	}

	return userRoutes
}
