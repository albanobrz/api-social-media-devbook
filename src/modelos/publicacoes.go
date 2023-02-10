package modelos

import (
	"errors"
	"strings"
	"time"
)

type Publicacao struct {
	ID        uint64    `json:"id,omitempty"`
	Titulo    string    `json:"titulo,omitempty"`
	Conteudo  string    `json:"conteudo,omitempty"`
	AutorID   uint64    `json:"autorId,omitempty"`
	AutorNick string    `json:"autorNick,omitempty"`
	Curtidas  uint64    `json:"curtidas"`
	CriadaEm  time.Time `json:"criadaEm,omitempty"`
}

func (publicacao *Publicacao) Preparar() error {
	if erro := publicacao.validar(); erro != nil {
		return erro
	}

	publicacao.formatar()
	return nil
}

func (Publicacao *Publicacao) validar() error {
	if Publicacao.Titulo == "" {
		return errors.New("O título é obrigatório e não pode estar em branco")
	}

	if Publicacao.Conteudo == "" {
		return errors.New("O contúdo é obrigatório e não pode estar em branco")
	}

	return nil
}

func (Publicacao *Publicacao) formatar() {
	Publicacao.Titulo = strings.TrimSpace(Publicacao.Titulo)
	Publicacao.Conteudo = strings.TrimSpace(Publicacao.Conteudo)
}
