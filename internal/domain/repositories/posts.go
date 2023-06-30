package repositories

import "api/internal/domain/entities"

type PostsRepository interface {
	Create(post entities.Post) (entities.Post, error)
	GetPosts(nick string) ([]entities.Post, error)
	GetPostWithId(id string) (entities.Post, error)
	UpdatePost(postID string, updatedPost entities.Post) error
	DeletePost(postID string) error
	GetAllPosts() ([]entities.Post, error)
	Like(postID string) error
	Dislike(postID string) error
}
