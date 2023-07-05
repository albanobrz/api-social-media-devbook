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
			URI:          "/mongo/users",
			Method:       http.MethodPost,
			Controller:   controllers.CreateUser,
			RequiresAuth: false,
		},
		{
			URI:          "/mongo/users/",
			Method:       http.MethodGet,
			Controller:   controllers.GetAllUsers,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}",
			Method:       http.MethodGet,
			Controller:   controllers.GetUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}",
			Method:       http.MethodPut,
			Controller:   controllers.UpdateUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}",
			Method:       http.MethodDelete,
			Controller:   controllers.DeleteUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/follow",
			Method:       http.MethodPost,
			Controller:   controllers.FollowUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/unfollow",
			Method:       http.MethodDelete,
			Controller:   controllers.UnfollowUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/followers",
			Method:       http.MethodGet,
			Controller:   controllers.GetFollowers,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/following",
			Method:       http.MethodGet,
			Controller:   controllers.GetFollowing,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/update-password",
			Method:       http.MethodPost,
			Controller:   controllers.UpdatePassword,
			RequiresAuth: true,
		},
	}

	return userRoutes
}
