package entities

import (
	"errors"
	"strings"
	"time"

	"api/internal/domain/security"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        uint64    `json:"id,omitempty" bson:"id"`
	Name      string    `json:"name,omitempty" bson:"name"`
	Nick      string    `json:"nick,omitempty" bson:"nick"`
	Email     string    `json:"email,omitempty" bson:"email"`
	Password  string    `json:"password,omitempty" bson:"password"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
	Followers []string  `json:"followers" bson:"followers"`
	Following []string  `json:"following" bson:"following"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	if err := user.format(step); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("The name is required and can't be empty")
	}

	if user.Nick == "" {
		return errors.New("The nick is required and can't be empty")
	}

	if user.Email == "" {

		return errors.New("The email is required and can't be empty")
	}
	if step == "createUser" && user.Password == "" {
		return errors.New("The password is required and can't be empty")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("The inserted email is invalid")
	}

	return nil
}

func (user *User) format(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if step == "createUser" {
		passwordHash, err := security.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(passwordHash)
	}
	return nil
}
