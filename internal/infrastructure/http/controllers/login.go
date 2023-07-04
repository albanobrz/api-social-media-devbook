package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/entities"
	"api/internal/domain/repositories"
	"api/internal/domain/security"
	"api/internal/infrastructure/http/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LoginController struct {
	userRepository repositories.UsersRepository
}

func NewLoginController(userRepository repositories.UsersRepository) *LoginController {
	return &LoginController{
		userRepository,
	}
}

func (controller *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}

	var user entities.User
	if err = json.Unmarshal(reqbody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userSavedOnDB, err := controller.userRepository.SearchByEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = security.VerifyPassword(userSavedOnDB.Password, user.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CreateTokenWithNick(userSavedOnDB.Nick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
