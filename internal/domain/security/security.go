package security

import "golang.org/x/crypto/bcrypt"

func Hash(hash string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(hash), bcrypt.DefaultCost)
}

func VerifyPassword(hashPassword, stringPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(stringPassword))
}
