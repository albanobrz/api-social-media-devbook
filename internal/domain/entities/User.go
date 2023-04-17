package entities

import (
	"errors"
	"strings"
	"time"

	"api/internal/domain/security"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuário utilizando a rede social
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// Prepare vai chamar os métodos para validar e formatar os usuário recebidos
func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	if err := user.format(step); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(etapa string) error {
	if user.Name == "" {
		return errors.New("The name is required and can't be empty")
	}

	if user.Nick == "" {
		return errors.New("The nick is required and can't be empty")
	}

	if user.Email == "" {

		return errors.New("The email is required and can't be empty")
	}
	if etapa == "cadastro" && user.Password == "" {
		return errors.New("The password is required and can't be empty")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("The inserted email is invalid")
	}

	return nil
}

func (user *User) format(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if step == "cadastro" {
		senhaHash, err := security.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(senhaHash)
	}
	return nil
}
