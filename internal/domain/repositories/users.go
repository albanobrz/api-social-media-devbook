package repositories

import "api/internal/domain/entities"

type UsersRepository interface {
	Create(user entities.User) (entities.User, error)
	SearchByEmail(email string) (entities.User, error)
	GetAllUsers() ([]entities.User, error)
	GetUserByNick(nick string) (entities.User, error)
	UpdateUser(nick string, user entities.User) error
	DeleteUser(nick string) error
	Follow(followerID string, followedID string) error
	Unfollow(unfollowerID string, unfollowedID string) error
	GetFollowers(userID string) ([]string, error)
	GetFollowing(userID string) ([]string, error)
	GetPassword(nick string) (string, error)
	UpdatePassword(nick string, password string) error
}
