package entities

import (
	"errors"
	"strings"
	"time"

	"api/src/security"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuário utilizando a rede social
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"nome,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"senha,omitempty"`
	CreatedAt time.Time `json:"CriadoEm,omitempty"`
}

// Prepare vai chamar os métodos para validar e formatar os usuário recebidos
func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	if erro := user.format(step); erro != nil {
		return erro
	}
	return nil
}

func (user *User) validate(etapa string) error {
	if user.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	if user.Nick == "" {
		return errors.New("O nick é obrigatório e não pode estar em branco")
	}

	if user.Email == "" {

		return errors.New("O email é obrigatório e não pode estar em branco")
	}
	if etapa == "cadastro" && user.Password == "" {
		return errors.New("A senha é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("O email inserido é inválido")
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
