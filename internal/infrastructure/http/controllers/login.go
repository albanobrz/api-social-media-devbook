package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/entities"
	"api/internal/domain/security"
	database "api/internal/infrastructure/database"
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}

	var user entities.User
	if err = json.Unmarshal(reqbody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}

	repository := repositories.NewUsersRepository(db)
	userSavedOnDB, err := repository.SearchByEmail(user.Email)
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
