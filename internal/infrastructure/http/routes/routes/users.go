package routes

import (
	"api/internal/infrastructure/http/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:          "/users",
		Method:       http.MethodPost,
		Function:     controllers.CreateUser,
		RequiresAuth: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodGet,
		Function:     controllers.GetUsers,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}",
		Method:       http.MethodGet,
		Function:     controllers.GetUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}",
		Method:       http.MethodPut,
		Function:     controllers.UpdateUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/follow",
		Method:       http.MethodPost,
		Function:     controllers.FollowUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/unfollow",
		Method:       http.MethodDelete,
		Function:     controllers.StopFollowingUser,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/followers",
		Method:       http.MethodGet,
		Function:     controllers.GetFollowers,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/following",
		Method:       http.MethodGet,
		Function:     controllers.GetFollowing,
		RequiresAuth: true,
	},
	{
		URI:          "/users/{userID}/update-password",
		Method:       http.MethodPost,
		Function:     controllers.UpdatePassword,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/users",
		Method:       http.MethodPost,
		Function:     controllers.CreateMongoUser,
		RequiresAuth: false,
	},
	{
		URI:          "/mongo/users/",
		Method:       http.MethodGet,
		Function:     controllers.GetAllUsersMongo,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/users/{userID}",
		Method:       http.MethodGet,
		Function:     controllers.GetUserMongo,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/users/{userID}",
		Method:       http.MethodPut,
		Function:     controllers.UpdateUserMongo,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/users/{userID}",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteUserMongo,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/users/{userID}/follow",
		Method:       http.MethodPost,
		Function:     controllers.FollowUserMongo,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/users/{userID}/unfollow",
		Method:       http.MethodDelete,
		Function:     controllers.UnfollowUserMongo,
		RequiresAuth: true,
	},
}
