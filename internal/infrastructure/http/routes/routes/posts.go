package routes

import (
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/controllers"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func ConfigPostsRoutes(db *mongo.Database) []Route {

	repository := repositories.NewPostsRepository(db)

	controllers := controllers.NewPostsController(*repository)

var PostsRoutes = []Route{
	{
		URI:          "/mongo/posts",
		Method:       http.MethodPost,
		Function:     controllers.CreatePost,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts/{userID}",
		Method:       http.MethodGet,
		Function:     controllers.GetPosts,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts/{postID}",
		Method:       http.MethodPut,
		Function:     controllers.UpdatePost,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts/{postID}",
		Method:       http.MethodDelete,
		Function:     controllers.DeletePost,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts/specific/{postID}",
		Method:       http.MethodGet,
		Function:     controllers.GetPost,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts",
		Method:       http.MethodGet,
		Function:     controllers.GetAllPosts,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts/{postID}/like",
		Method:       http.MethodPost,
		Function:     controllers.LikePost,
		RequiresAuth: true,
	},
	{
		URI:          "/mongo/posts/{postID}/dislike",
		Method:       http.MethodPost,
		Function:     controllers.DislikePost,
		RequiresAuth: true,
	},
}
