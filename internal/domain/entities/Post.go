package entities

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID         uint64    `json:"id,omitempty" bson:"id"`
	Title      string    `json:"title,omitempty" bson:"title"`
	Content    string    `json:"content,omitempty" bson:"content"`
	AuthorID   string    `json:"authorId,omitempty" bson:"authorId"`
	AuthorNick string    `json:"authorNick,omitempty" bson:"authorNick"`
	Likes      uint64    `json:"likes" bson:"likes"`
	CreatedAt  time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
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
