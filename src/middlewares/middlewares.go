package middlewares

import (
	"api/internal/infrastructure/http/responses"
	"api/src/auth"
	"log"
	"net/http"
)

// Logger escreve informações da requisição no terminal
func Logger(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		proximaFuncao(w, r)
	}
}

// Verifica se o usuário fazendo a requisição está autenticado
func Authenticate(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	// Handlerfunc é o (w, r) comum das outras requisições
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := auth.ValidateToken(r); erro != nil {
			responses.Error(w, http.StatusUnauthorized, erro)
			return
		}
		proximaFuncao(w, r)
	}
}
