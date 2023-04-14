package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `bson:"name" json:"name"`
	ID               uint64 `bson:"id" json:"id"`
	Nick             string `bson:"nick" json:"nick"`
	Email            string `bson:"email" json:"email"`
	Password         string `bson:"password" json:"password"`
	CreatedAt        string `bson:"createdAt" json:"createdAt"`
}

func NewUser(name string, nick string) *User {
	return &User{
		Name: name,
		Nick: nick,
	}
}
