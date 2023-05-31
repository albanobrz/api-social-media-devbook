package repositories

import (
	"api/internal/domain/entities"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type PostsRepository struct {
	collection *mongo.Collection
}

func NewPostsRepository(db *mongo.Database) *PostsRepository {
	collection := db.Collection("posts")
	return &PostsRepository{
		collection,
	}
}

func (repository *PostsRepository) Create(post entities.Post) (entities.Post, error) {
	newPost := entities.Post{
		Title:      post.Title,
		Content:    post.Content,
		AuthorID:   post.AuthorID,
		AuthorNick: post.AuthorID,
		Likes:      0,
		CreatedAt:  time.Now(),
	}

	_, err := repository.collection.InsertOne(context.Background(), newPost)
	if err != nil {
		return entities.Post{}, err
	}

	return newPost, nil
}

func (repository *PostsRepository) GetPosts(nick string) ([]entities.Post, error) {
	cursor, err := repository.collection.Find(context.TODO(), bson.M{"authorNick": nick})
	if err != nil {
		return []entities.Post{}, err
	}
	defer cursor.Close(context.TODO())

	var results []entities.Post
	for cursor.Next(context.TODO()) {
		var result entities.Post
		err := cursor.Decode(&result)
		if err != nil {
			return []entities.Post{}, err
		}

		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return []entities.Post{}, err
	}

	return results, nil
}

func (repository *PostsRepository) GetPostWithId(id string) (entities.Post, error) {
	idString, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entities.Post{}, err
	}

	var result entities.Post
	err = repository.collection.FindOne(context.TODO(), bson.M{"_id": idString}).Decode(&result)
	if err != nil {
		return entities.Post{}, fmt.Errorf("This post doens't exists")
	}

	return result, nil
}

func (repository *PostsRepository) UpdatePost(postID string, updatedPost entities.Post) error {
	idString, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"title": updatedPost.Title, "content": updatedPost.Content, "updatedAt": time.Now()}}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"_id": idString}, update)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostsRepository) DeletePost(postID string) error {
	idString, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	_, err = repository.collection.DeleteOne(context.TODO(), bson.M{"_id": idString})
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostsRepository) GetAllPosts() ([]entities.Post, error) {
	cursor, err := repository.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []entities.Post

	for cursor.Next(context.Background()) {
		var post entities.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repository *PostsRepository) Like(postID string) error {
	idString, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"_id": idString}, bson.M{"$inc": bson.M{"likes": 1}})
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostsRepository) Dislike(postID string) error {
	idString, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	var result entities.Post
	err = repository.collection.FindOne(context.TODO(), bson.M{"_id": idString}).Decode(&result)
	if err != nil {
		return err
	}

	if result.Likes <= 0 {
		return fmt.Errorf("Like count is already 0")
	}

	update := bson.M{
		"$inc": bson.M{"likes": -1},
	}

	_, err = repository.collection.UpdateOne(context.TODO(), bson.M{"_id": idString}, update)
	if err != nil {
		return err
	}

	return nil
}
