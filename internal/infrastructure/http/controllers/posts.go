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

// func CreatePost(w http.ResponseWriter, r *http.Request) {
// 	userID, err := auth.GetUserID(r)
// 	if err != nil {
// 		responses.Error(w, http.StatusUnauthorized, err)
// 		return
// 	}

// 	reqbody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		responses.Error(w, http.StatusUnprocessableEntity, err)
// 		return
// 	}

// 	var post entities.Post
// 	if err = json.Unmarshal(reqbody, &post); err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	post.AuthorID = userID

// 	if err = post.Prepare(); err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	post.ID, err = repository.Create(post)
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusCreated, post)
// }

// func GetPosts(w http.ResponseWriter, r *http.Request) {
// 	userID, err := auth.GetUserID(r)
// 	if err != nil {
// 		responses.Error(w, http.StatusUnauthorized, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	posts, err := repository.Get(userID)
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusOK, posts)
// }

// func GetPost(w http.ResponseWriter, r *http.Request) {
// 	// o params e a string é como passou na rota
// 	params := mux.Vars(r)
// 	postID, err := strconv.ParseUint(params["postID"], 10, 64)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	post, err := repository.GetWithID(postID)
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusOK, post)
// }

// func UpdatePost(w http.ResponseWriter, r *http.Request) {
// 	userID, err := auth.GetUserID(r)
// 	if err != nil {
// 		responses.Error(w, http.StatusUnauthorized, err)
// 		return
// 	}

// 	params := mux.Vars(r)
// 	postID, err := strconv.ParseUint(params["postID"], 10, 64)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	postSavedOnDB, err := repository.GetWithID(postID)
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	if postSavedOnDB.AuthorID != userID {
// 		responses.Error(w, http.StatusForbidden, errors.New("It's not possible to update other's posts"))
// 		return
// 	}

// 	reqbody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		responses.Error(w, http.StatusUnprocessableEntity, err)
// 		return
// 	}

// 	var post entities.Post
// 	if err = json.Unmarshal(reqbody, &post); err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if err = post.Prepare(); err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if err = repository.Update(postID, post); err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusNoContent, nil)
// }

// func DeletePost(w http.ResponseWriter, r *http.Request) {
// 	userID, err := auth.GetUserID(r)
// 	if err != nil {
// 		responses.Error(w, http.StatusUnauthorized, err)
// 		return
// 	}

// 	params := mux.Vars(r)
// 	postID, err := strconv.ParseUint(params["postID"], 10, 64)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	postSavedOnDB, err := repository.GetWithID(postID)
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	if postSavedOnDB.AuthorID != userID {
// 		responses.Error(w, http.StatusForbidden, errors.New("It's not possible deleting other's post"))
// 		return
// 	}

// 	if err = repository.Delete(postID); err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusNoContent, nil)
// }

// func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	userID, err := strconv.ParseUint(params["userID"], 10, 64)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	posts, err := repository.GetByUser(userID)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusOK, posts)
// }

// func LikePost(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	postID, err := strconv.ParseUint(params["postID"], 10, 64)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	if err = repository.Like(postID); err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusNoContent, nil)
// }

// func DislikePost(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	postID, err := strconv.ParseUint(params["postID"], 10, 64)
// 	if err != nil {
// 		responses.Error(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	db, err := database.Connect()
// 	if err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	repository := repositories.NewPostsRepository(db)
// 	if err = repository.Dislike(postID); err != nil {
// 		responses.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	responses.JSON(w, http.StatusNoContent, nil)
// }

func CreatePostMongo(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	post, err = repository.CreateMongo(post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func GetPostsMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userNick := params["userID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	posts, err := repository.GetPostsMongo(userNick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func UpdatePostMongo(w http.ResponseWriter, r *http.Request) {
	userNick, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	postSavedOnDB, err := repository.GetPostWithIdMongo(postID)
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

	if err = repository.UpdatePostMongo(postID, post); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeletePostMongo(w http.ResponseWriter, r *http.Request) {
	userNick, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	postSavedOnDB, err := repository.GetPostWithIdMongo(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	if postSavedOnDB.AuthorNick != userNick {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible deleting other's post"))
		return
	}

	if err = repository.DeletePostMongo(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetPostMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	post, err := repository.GetPostWithIdMongo(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func GetAllPostsMongo(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	posts, err := repository.GetAllPostsMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func LikePostMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	err = repository.Like(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}

func DislikePostMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewPostsRepositoryMongo(db)
	err = repository.Dislike(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, nil)
}
