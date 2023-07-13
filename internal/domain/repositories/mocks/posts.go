package mocks

import (
	"api/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

type PostsRepositoryMock struct {
	mock.Mock
}

func NewPostsRepositoryMock() *PostsRepositoryMock {
	return &PostsRepositoryMock{}
}

func (repository *PostsRepositoryMock) CreatePost(post entities.Post) (entities.Post, error) {
	args := repository.Called(post)
	return args.Get(0).(entities.Post), args.Error(1)
}

func (repository *PostsRepositoryMock) GetPosts(nick string) ([]entities.Post, error) {
	args := repository.Called(nick)
	return args.Get(0).([]entities.Post), args.Error(1)
}

func (repository *PostsRepositoryMock) GetPostWithId(id string) (entities.Post, error) {
	args := repository.Called(id)
	return args.Get(0).(entities.Post), args.Error(1)
}

func (repository *PostsRepositoryMock) UpdatePost(postID string, updatedPost entities.Post) error {
	args := repository.Called(postID, updatedPost)
	return args.Error(0)
}

func (repository *PostsRepositoryMock) DeletePost(postID string) error {
	args := repository.Called(postID)
	return args.Error(0)
}

func (repository *PostsRepositoryMock) GetAllPosts() ([]entities.Post, error) {
	args := repository.Called()
	return args.Get(0).([]entities.Post), args.Error(1)
}

func (repository *PostsRepositoryMock) Like(postID string) error {
	args := repository.Called(postID)
	return args.Error(0)
}

func (repository *PostsRepositoryMock) Dislike(postID string) error {
	args := repository.Called(postID)
	return args.Error(0)
}
