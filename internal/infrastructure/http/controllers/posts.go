package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/entities"
	"api/internal/domain/repositories"
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

	post, err = controller.PostRepository.CreatePost(post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func (controller *PostsController) GetPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userNick := params["userID"]

	posts, err := controller.PostRepository.GetPosts(userNick)
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

	postSavedOnDB, err := controller.PostRepository.GetPostWithId(postID)
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

	if err = controller.PostRepository.UpdatePost(postID, post); err != nil {
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

	postSavedOnDB, err := controller.PostRepository.GetPostWithId(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	if postSavedOnDB.AuthorNick != userNick {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible deleting other's post"))
		return
	}

	if err = controller.PostRepository.DeletePost(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *PostsController) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	post, err := controller.PostRepository.GetPostWithId(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func (controller *PostsController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := controller.PostRepository.GetAllPosts()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (controller *PostsController) LikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	err := controller.PostRepository.Like(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}

func (controller *PostsController) DislikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	err := controller.PostRepository.Dislike(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}
