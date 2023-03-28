package routes

import (
	"api/internal/infrastructure/http/controllers"
	"net/http"
)

var PostsRoutes = []Route{
	{
		URI:          "/publicacoes",
		Method:       http.MethodPost,
		Function:     controllers.CreatePost,
		RequiresAuth: true,
	},
	{
		URI:          "/publicacoes",
		Method:       http.MethodGet,
		Function:     controllers.GetPosts,
		RequiresAuth: true,
	},
	{
		URI:          "/publicacoes/{publicacaoId}",
		Method:       http.MethodGet,
		Function:     controllers.GetPost,
		RequiresAuth: true,
	},
	{
		URI:          "/publicacoes/{publicacaoId}",
		Method:       http.MethodPut,
		Function:     controllers.UpdatePost,
		RequiresAuth: true,
	},
	{
		URI:          "/publicacoes/{publicacaoId}",
		Method:       http.MethodDelete,
		Function:     controllers.DeletePost,
		RequiresAuth: true,
	},
	{
		URI:          "/usuarios/{usuarioId}/publicacoes",
		Method:       http.MethodGet,
		Function:     controllers.GetPostsByUser,
		RequiresAuth: true,
	},
	{
		URI:          "/publicacoes/{publicacaoId}/curtir",
		Method:       http.MethodPost,
		Function:     controllers.LikePost,
		RequiresAuth: true,
	},
	{
		URI:          "/publicacoes/{publicacaoId}/descurtir",
		Method:       http.MethodPost,
		Function:     controllers.DislikePost,
		RequiresAuth: true,
	},
}
