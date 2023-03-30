package entities

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID         uint64    `json:"id,omitempty"`
	Title      string    `json:"titulo,omitempty"`
	Content    string    `json:"conteudo,omitempty"`
	AuthorID   uint64    `json:"autorId,omitempty"`
	AuthorNick string    `json:"autorNick,omitempty"`
	Likes      uint64    `json:"curtidas"`
	CreatedAt  time.Time `json:"criadaEm,omitempty"`
}

func (post *Post) Prepare() error {
	if err := post.validate(); err != nil {
		return err
	}

	post.format()
	return nil
}

func (post *Post) validate() error {
	if post.Title == "" {
		return errors.New("The title is required and can't be empty")
	}

	if post.Content == "" {
		return errors.New("The content is required and can't be empty")
	}

	return nil
}

func (post *Post) format() {
	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)
}
