package repositories

import (
	"api/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type UsersRepository struct {
	collection *mongo.Collection
}

func NewUsersRepository(db *mongo.Database) *UsersRepository {
	collection := db.Collection("users")
	return &UsersRepository{
		collection,
	}
}

func (repository *UsersRepository) Create(user entities.User) (entities.User, error) {
	existingUser := entities.User{}
	err := repository.collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return entities.User{}, errors.New("there is already a user with the same email")
	}

	err = repository.collection.FindOne(context.Background(), bson.M{"nick": user.Nick}).Decode(&existingUser)
	if err == nil {
		return entities.User{}, errors.New("there is already a user with the same nick")
	}

	newUser := entities.User{
		Name:      user.Name,
		Nick:      user.Nick,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		Followers: []string{},
		Following: []string{},
	}

	_, err = repository.collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return entities.User{}, err
	}

	return newUser, nil
}

func (repository *UsersRepository) SearchByEmail(email string) (entities.User, error) {
	existingUser := entities.User{}

	err := repository.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		return entities.User{}, err
	}

	return existingUser, nil
}

func (repository *UsersRepository) GetAllUsers() ([]entities.User, error) {
	cursor, err := repository.collection.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{"password": 0}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []entities.User

	for cursor.Next(context.Background()) {
		var user entities.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *UsersRepository) GetUserByNick(nick string) (entities.User, error) {
	existingUser := entities.User{}

	options := options.FindOne().SetProjection(bson.M{"password": 0})

	err := repository.collection.FindOne(context.Background(), bson.M{"nick": nick}, options).Decode(&existingUser)
	if err != nil {
		return entities.User{}, err
	}

	return existingUser, nil
}

func (repository *UsersRepository) UpdateUser(nick string, user entities.User) error {
	update := bson.M{
		"$set": bson.M{"name": user.Name, "nick": user.Nick, "email": user.Email, "updatedAt": time.Now()},
	}

	_, err := repository.collection.UpdateOne(context.Background(), bson.M{"nick": nick}, update)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UsersRepository) DeleteUser(nick string) error {
	_, err := repository.collection.DeleteOne(context.Background(), bson.M{"nick": nick})
	if err != nil {
		return err
	}

	return nil
}

func (repository *UsersRepository) Follow(followerID string, followedID string) error {
	filter := bson.M{
		"nick":      followerID,
		"following": bson.M{"$in": []string{followedID}},
	}

	var result entities.User
	err := repository.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		return fmt.Errorf("You already follow this user")
	}

	// update follower data
	updateFollower := bson.M{
		"$push": bson.M{
			"following": followedID,
		},
	}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"nick": followerID}, updateFollower)
	if err != nil {
		return err
	}

	// update followed data
	updateFollowed := bson.M{
		"$push": bson.M{
			"followers": followerID,
		},
	}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"nick": followedID}, updateFollowed)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UsersRepository) Unfollow(unfollowerID string, unfollowedID string) error {
	filter := bson.M{
		"nick":      unfollowerID,
		"following": bson.M{"$in": []string{unfollowedID}},
	}

	var result entities.User
	err := repository.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return fmt.Errorf("You don't follow this user already")
	}

	// update follower data
	updateFollower := bson.M{
		"$pull": bson.M{
			"following": unfollowedID,
		},
	}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"nick": unfollowerID}, updateFollower)
	if err != nil {
		return err
	}

	// update followed data
	updateFollowed := bson.M{
		"$pull": bson.M{
			"followers": unfollowerID,
		},
	}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"nick": unfollowedID}, updateFollowed)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UsersRepository) GetFollowers(userID string) ([]string, error) {
	var followers []string
	cursor, err := repository.collection.Find(context.TODO(), bson.M{"nick": userID})
	if err != nil {
		return []string{}, err
	}

	for cursor.Next(context.TODO()) {
		var result entities.User
		err := cursor.Decode(&result)
		if err != nil {
			return []string{}, err
		}

		followers = append(followers, result.Followers...)
	}

	if err := cursor.Err(); err != nil {
		return []string{}, err
	}

	cursor.Close(context.TODO())

	return followers, nil
}

func (repository *UsersRepository) GetFollowing(userID string) ([]string, error) {

	var following []string
	cursor, err := repository.collection.Find(context.TODO(), bson.M{"nick": userID})
	if err != nil {
		return []string{}, err
	}

	for cursor.Next(context.TODO()) {
		var result entities.User
		err := cursor.Decode(&result)
		if err != nil {
			return []string{}, err
		}

		following = append(following, result.Following...)
	}

	if err := cursor.Err(); err != nil {
		return []string{}, err
	}

	cursor.Close(context.TODO())

	return following, nil
}

func (repository *UsersRepository) GetPassword(nick string) (string, error) {
	existingUser := entities.User{}

	err := repository.collection.FindOne(context.Background(), bson.M{"nick": nick}).Decode(&existingUser)
	if err != nil {
		return "", err
	}

	return existingUser.Password, nil
}

func (repository *UsersRepository) UpdatePassword(nick string, password string) error {
	updatePassword := bson.M{
		"$set": bson.M{
			"password":  password,
			"updatedAt": time.Now(),
		},
	}

	_, err := repository.collection.UpdateOne(context.TODO(), bson.M{"nick": nick}, updatePassword)
	if err != nil {
		return err
	}

	return nil
}
