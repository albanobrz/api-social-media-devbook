package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/entities"
	"api/internal/domain/security"
	database "api/internal/infrastructure/database"
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user entities.User
	if err = json.Unmarshal(reqbody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := user.Prepare("createUser"); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	newUser, err := repository.Create(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newUser)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	repository := repositories.NewUsersRepository(db)
	users, err := repository.GetAllUsers()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	nick := params["userID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	repository := repositories.NewUsersRepository(db)
	user, err := repository.GetUserByNick(nick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nick := params["userID"]

	userNickOnToken, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	fmt.Println(nick, userNickOnToken)

	if nick != userNickOnToken {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible update another user"))
		return
	}

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user entities.User
	if err = json.Unmarshal(reqbody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("edicao"); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if err = repository.UpdateUser(nick, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	userNickOnToken, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userNickOnToken {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible delete another user"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if err = repository.DeleteUser(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	followedID := params["userID"]

	if followerID == followedID {
		responses.Error(w, http.StatusForbidden, errors.New("You can't follow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if err = repository.Follow(followerID, followedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	unfollowerID, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	unfollowedID := params["userID"]

	if unfollowerID == unfollowedID {
		responses.Error(w, http.StatusForbidden, errors.New("You can't unfollow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	if err = repository.Unfollow(unfollowerID, unfollowedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	followers, err := repository.GetFollowers(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	following, err := repository.GetFollowing(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, following)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userNickOnToken, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userNick := params["userID"]
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userNickOnToken != userNick {
		responses.Error(w, http.StatusForbidden, errors.New("You can't update passwords from another user"))
		return
	}

	reqbody, err := ioutil.ReadAll(r.Body)

	var password entities.Password
	if err = json.Unmarshal(reqbody, &password); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepository(db)
	passwordSavedOnDB, err := repository.GetPassword(userNick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassword(passwordSavedOnDB, password.Current); err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("The current password doesn't match"))
		return
	}

	HashedPassword, err := security.Hash(password.New)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.UpdatePassword(userNick, string(HashedPassword)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
