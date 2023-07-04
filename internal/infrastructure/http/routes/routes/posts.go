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
			Controller:   controllers.CreatePost,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts/{userID}",
			Method:       http.MethodGet,
			Controller:   controllers.GetPosts,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts/{postID}",
			Method:       http.MethodPut,
			Controller:   controllers.UpdatePost,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts/{postID}",
			Method:       http.MethodDelete,
			Controller:   controllers.DeletePost,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts/specific/{postID}",
			Method:       http.MethodGet,
			Controller:   controllers.GetPost,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts",
			Method:       http.MethodGet,
			Controller:   controllers.GetAllPosts,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts/{postID}/like",
			Method:       http.MethodPost,
			Controller:   controllers.LikePost,
			RequiresAuth: true,
		},
		{
			URI:          "/mongo/posts/{postID}/dislike",
			Method:       http.MethodPost,
			Controller:   controllers.DislikePost,
			RequiresAuth: true,
		},
	}
	return PostsRoutes
}
