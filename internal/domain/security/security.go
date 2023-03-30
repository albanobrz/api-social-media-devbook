package security

import "golang.org/x/crypto/bcrypt"

// Hash recebe uma string e coloca um hash nela
func Hash(hash string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(hash), bcrypt.DefaultCost)
}

// Verificar senha compara uma senha e um hash e retorna se elas são iguais
func VerifyPassword(hashPassword, stringPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(stringPassword))
}
