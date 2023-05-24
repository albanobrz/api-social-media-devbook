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
	"strconv"
	"strings"

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
	user.ID, err = repository.Create(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	users, err := repository.Get(nameOrNick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	user, err := repository.GetByID(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDOnToken, err := auth.GetUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDOnToken {
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
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err = repository.Update(userID, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDOnToken, err := auth.GetUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDOnToken {
		responses.Error(w, http.StatusForbidden, errors.New("It's not possible delete another user"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err = repository.Delete(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	// quem vai seguir é quem ta fazendo a requisição (tem o id), e quem vai ser seguido, é o usuário que o ID ta no parametro
	followedID, err := auth.GetUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if followedID == userID {
		responses.Error(w, http.StatusForbidden, errors.New("You can't follow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err = repository.Follow(userID, followedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func StopFollowingUser(w http.ResponseWriter, r *http.Request) {
	followedID, err := auth.GetUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if followedID == userID {
		responses.Error(w, http.StatusForbidden, errors.New("You can't unfollow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err = repository.StopFollow(userID, followedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

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
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	following, err := repository.GetFollowing(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, following)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userIDOnToken, err := auth.GetUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userIDOnToken != userID {
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
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	passwordSavedOnDB, err := repository.GetPassword(userID)
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

	if err = repository.UpdatePassword(userID, string(HashedPassword)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func CreateMongoUser(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	newUser, err := repository.CreateMongo(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newUser)
}

func GetAllUsersMongo(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	users, err := repository.GetAllUsersMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, users)
}

func GetUserMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	nick := params["userID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	user, err := repository.GetUserByNick(nick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	responses.JSON(w, http.StatusOK, user)
}

func UpdateUserMongo(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	if err = repository.UpdateUserMongo(nick, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeleteUserMongo(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	if err = repository.DeleteUserMongo(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUserMongo(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	if err = repository.FollowMongo(followerID, followedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnfollowUserMongo(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	if err = repository.UnfollowMongo(unfollowerID, unfollowedID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetFollowersMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	followers, err := repository.GetFollowersMongo(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func GetFollowingMongo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	db, err := database.ConnectMongo()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	repository := repositories.NewUsersRepositoryMongo(db)
	following, err := repository.GetFollowingMongo(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, following)
}
