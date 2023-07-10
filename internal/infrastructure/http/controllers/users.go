package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/entities"
	"api/internal/domain/repositories"
	"api/internal/domain/security"
	"api/internal/infrastructure/http/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type UsersController struct {
	userRepository repositories.UsersRepository
}

func NewUsersController(userRepository repositories.UsersRepository) *UsersController {
	return &UsersController{
		userRepository,
	}
}

func (controller *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	newUser, err := controller.userRepository.Create(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newUser)
}

func (controller *UsersController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := controller.userRepository.GetAllUsers()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, users)
}

func (controller *UsersController) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	nick := params["userID"]

	user, err := controller.userRepository.GetUserByNick(nick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, user)
}

func (controller *UsersController) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	if err = controller.userRepository.UpdateUser(nick, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	if err = controller.userRepository.DeleteUser(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *UsersController) FollowUser(w http.ResponseWriter, r *http.Request) {
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

	if err = controller.userRepository.Follow(followerID, followedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *UsersController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
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

	if err = controller.userRepository.Unfollow(unfollowerID, unfollowedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller *UsersController) GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	if userID == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("You must insert the user ID"))
		return
	}

	followers, err := controller.userRepository.GetFollowers(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func (controller *UsersController) GetFollowing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	if userID == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("You must insert the user ID"))
		return
	}

	following, err := controller.userRepository.GetFollowing(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, following)
}

func (controller *UsersController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userNickOnToken, err := auth.GetUserNick(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userNick := params["userID"]
	if userNick == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("userID is missing"))
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

	passwordSavedOnDB, err := controller.userRepository.GetPassword(userNick)
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

	if err = controller.userRepository.UpdatePassword(userNick, string(HashedPassword)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
