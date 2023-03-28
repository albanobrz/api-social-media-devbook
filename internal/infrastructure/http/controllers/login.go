package controllers

import (
	"api/internal/domain/entities"
	database "api/internal/infrastructure/database"
	"api/internal/infrastructure/database/repositories"
	"api/internal/infrastructure/http/responses"
	"api/src/auth"
	"api/src/security"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Login é responsável por autenticar um usuário na API
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
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	userSavedOnDB, err := repository.BuscarPorEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = security.VerifyPassword(userSavedOnDB.Password, user.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CreateToken(userSavedOnDB.ID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))

	// Vem um hash quando faz a requisição. Tem que comparar o hash da requisição com o do banco de dados pra validar
}
