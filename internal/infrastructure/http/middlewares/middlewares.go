package middlewares

import (
	"api/internal/application/auth"
	"api/internal/infrastructure/http/responses"
	"log"
	"net/http"
)

// Logger escreve informações da requisição no terminal
func Logger(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		nextFunction(w, r)
	}
}

// Verifica se o usuário fazendo a requisição está autenticado
func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	// Handlerfunc é o (w, r) comum das outras requisições
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}
		nextFunction(w, r)
	}
}
