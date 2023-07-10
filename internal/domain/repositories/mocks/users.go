package mocks

import (
	"api/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

type UsersRepositoryMock struct {
	mock.Mock
}

func NewUsersRepositoryMock() *UsersRepositoryMock {
	return &UsersRepositoryMock{}
}

func (repository *UsersRepositoryMock) Create(user entities.User) (entities.User, error) {
	args := repository.Called(user)
	return args.Get(0).(entities.User), args.Error(1)
}

func (repository *UsersRepositoryMock) SearchByEmail(email string) (entities.User, error) {
	args := repository.Called(email)
	return args.Get(0).(entities.User), args.Error(1)
}

func (repository *UsersRepositoryMock) GetAllUsers() ([]entities.User, error) {
	args := repository.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

func (repository *UsersRepositoryMock) GetUserByNick(nick string) (entities.User, error) {
	args := repository.Called(nick)
	return args.Get(0).(entities.User), args.Error(1)
}

func (repository *UsersRepositoryMock) UpdateUser(nick string, user entities.User) error {
	args := repository.Called(nick, user)
	return args.Error(0)
}

func (repository *UsersRepositoryMock) DeleteUser(nick string) error {
	args := repository.Called(nick)
	return args.Error(0)
}

func (repository *UsersRepositoryMock) Follow(followerID string, followedID string) error {
	args := repository.Called(followedID, followerID)
	return args.Error(0)
}

func (repository *UsersRepositoryMock) Unfollow(unfollowerID string, unfollowedID string) error {
	args := repository.Called(unfollowerID, unfollowedID)
	return args.Error(0)
}

func (repository *UsersRepositoryMock) GetFollowers(userID string) ([]string, error) {
	args := repository.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (repository *UsersRepositoryMock) GetFollowing(userID string) ([]string, error) {
	args := repository.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (repository *UsersRepositoryMock) GetPassword(nick string) (string, error) {
	args := repository.Called(nick)
	return args.Get(0).(string), args.Error(1)
}

func (repository *UsersRepositoryMock) UpdatePassword(nick string, password string) error {
	args := repository.Called(nick, password)
	return args.Error(0)
}

func (repository *UsersRepositoryMock) SecurityMock(passwordSavedOnDb string, newPassword string) error {
	args := repository.Called(passwordSavedOnDb, newPassword)
	return args.Error(0)
}
