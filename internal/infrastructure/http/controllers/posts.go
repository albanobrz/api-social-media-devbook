package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/entities"
	database "api/internal/infrastructure/database"
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/responses"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type PostsController struct {
	PostRepository repositories.PostsRepository
}

func NewPostsController(postRepository repositories.PostsRepository) *PostsController {
	return &PostsController{
		postRepository,
	}
}

func (controller *PostsController) CreatePost(w http.ResponseWriter, r *http.Request) {
	userNick, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post entities.Post
	if err = json.Unmarshal(reqbody, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorID = userNick

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	post, err = repository.CreatePost(post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func (controller *PostsController) GetPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userNick := params["userID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	posts, err := repository.GetPosts(userNick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (controller *PostsController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	userNick, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	postSavedOnDB, err := repository.GetPostWithId(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	if postSavedOnDB.AuthorNick != userNick {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible to update other's posts"))
		return
	}

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post entities.Post
	if err = json.Unmarshal(reqbody, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.UpdatePost(postID, post); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *PostsController) DeletePost(w http.ResponseWriter, r *http.Request) {
	userNick, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	postSavedOnDB, err := repository.GetPostWithId(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	if postSavedOnDB.AuthorNick != userNick {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible deleting other's post"))
		return
	}

	if err = repository.DeletePost(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *PostsController) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	post, err := repository.GetPostWithId(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func (controller *PostsController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	posts, err := repository.GetAllPosts()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (controller *PostsController) LikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	err = repository.Like(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}

func (controller *PostsController) DislikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepository(db)
	err = repository.Dislike(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}
