package routes

import (
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/controllers"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigUsersRoutes(db *mongo.Database) []Route {

	repository := repositories.NewUsersRepository(db)

	controllers := controllers.NewUsersController(*repository)

	var userRoutes = []Route{
		{
			URI:          "/mongo/users",
			Method:       http.MethodPost,
			Function:     controllers.CreateUser,
			RequiresAuth: false,
		},
		{
			URI:          "/mongo/users/",
			Method:       http.MethodGet,
			Function:     controllers.GetAllUsers,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}",
			Method:       http.MethodGet,
			Function:     controllers.GetUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}",
			Method:       http.MethodPut,
			Function:     controllers.UpdateUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}",
			Method:       http.MethodDelete,
			Function:     controllers.DeleteUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/follow",
			Method:       http.MethodPost,
			Function:     controllers.FollowUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/unfollow",
			Method:       http.MethodDelete,
			Function:     controllers.UnfollowUser,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/followers",
			Method:       http.MethodGet,
			Function:     controllers.GetFollowers,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/following",
			Method:       http.MethodGet,
			Function:     controllers.GetFollowing,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/users/{userID}/update-password",
			Method:       http.MethodPost,
			Function:     controllers.UpdatePassword,
			RequiresAuth: true,
		},
	}

	return userRoutes
}
