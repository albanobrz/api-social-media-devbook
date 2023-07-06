package controllers

import (
	"api/internal/application/auth"
	"api/internal/domain/security"
	"os"
	"testing"
)

var (
	ValidToken, DiffToken string
)

var Hashed string

func TestMain(m *testing.M) {
	ValidToken, _ = auth.CreateTokenWithNick("1")
	DiffToken, _ = auth.CreateTokenWithNick("2")
	hash, _ := security.Hash("123456")

	Hashed = string(hash)

	os.Exit(m.Run())
}
