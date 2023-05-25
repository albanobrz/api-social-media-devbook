package repositories

import (
	"api/internal/domain/entities"
	"context"
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Posts struct {
	db *sql.DB
}

func NewPostsRepository(db *sql.DB) *Posts {
	return &Posts{db}
}

type PostsRepositoryMongo struct {
	collection *mongo.Collection
}

func NewPostsRepositoryMongo(db *mongo.Database) *PostsRepositoryMongo {
	collection := db.Collection("posts")
	return &PostsRepositoryMongo{
		collection,
	}
}

// func (repository Posts) Create(post entities.Post) (uint64, error) {
// 	statement, err := repository.db.Prepare("insert into posts (title, content, author_id) values (?, ?, ?)")
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer statement.Close()

// 	result, err := statement.Exec(post.Title, post.Content, post.AuthorID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	lastInsertID, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return uint64(lastInsertID), nil
// }

// func (repository Posts) GetWithID(postID uint64) (entities.Post, error) {
// 	row, err := repository.db.Query(`
// 		select p.*, u.nick from
// 		posts p inner join users u
// 		on u.id = p.author_id where p.id = ?
// 	`, postID)
// 	if err != nil {
// 		return entities.Post{}, err
// 	}

// 	defer row.Close()

// 	var post entities.Post

// 	if row.Next() {
// 		if err = row.Scan(
// 			&post.ID,
// 			&post.Title,
// 			&post.Content,
// 			&post.AuthorID,
// 			&post.Likes,
// 			&post.CreatedAt,
// 			&post.AuthorNick,
// 		); err != nil {
// 			return entities.Post{}, err
// 		}
// 	}

// 	return post, nil
// }

// func (repository Posts) Get(userID uint64) ([]entities.Post, error) {
// 	rows, err := repository.db.Query(`
// 		select distinct p.*, u.nick from posts p
// 		inner join users u on u.id = p.author_id
// 		inner join followers s on p.author_id = s.user_id
// 		where u.id = ? or s.follower_id = ?
// 		order by 1 desc
// 	`, userID, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []entities.Post

// 	for rows.Next() {
// 		var post entities.Post
// 		if err = rows.Scan(
// 			&post.ID,
// 			&post.Title,
// 			&post.Content,
// 			&post.AuthorID,
// 			&post.Likes,
// 			&post.CreatedAt,
// 			&post.AuthorNick,
// 		); err != nil {
// 			return nil, err
// 		}

// 		posts = append(posts, post)
// 	}
// 	return posts, nil
// }

// func (repository Posts) Update(postID uint64, post entities.Post) error {
// 	statement, err := repository.db.Prepare("update posts set title = ?, content = ? where id = ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	if _, err = statement.Exec(post.Title, post.Content, postID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (repository Posts) Delete(postID uint64) error {
// 	statement, err := repository.db.Query(
// 		"delete from posts where id = ?", postID,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	return nil
// }

// func (repository Posts) GetByUser(userID uint64) ([]entities.Post, error) {
// 	rows, err := repository.db.Query(`
// 		select p.*, u.nick from posts p
// 		join users u on u.id = p.author_id
// 		where p.author_id = ?
// 	`, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []entities.Post

// 	for rows.Next() {
// 		var post entities.Post
// 		if err = rows.Scan(
// 			&post.ID,
// 			&post.Title,
// 			&post.Content,
// 			&post.AuthorID,
// 			&post.Likes,
// 			&post.CreatedAt,
// 			&post.AuthorNick,
// 		); err != nil {
// 			return nil, err
// 		}

// 		posts = append(posts, post)
// 	}
// 	return posts, nil
// }

// func (repository Posts) Like(postID uint64) error {
// 	statement, err := repository.db.Prepare("update posts set likes = likes + 1 where id = ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	if _, err = statement.Exec(postID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (repository Posts) Dislike(postID uint64) error {
// 	statement, err := repository.db.Prepare(`
// 		update posts set likes =
// 		CASE
// 			WHEN likes > 0 THEN likes - 1
// 			ELSE 0
// 		END
// 		where id = ?
// 	 `)
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	if _, err = statement.Exec(postID); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (repository *PostsRepositoryMongo) CreateMongo(post entities.Post) (entities.Post, error) {
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

func (repository *PostsRepositoryMongo) GetPostsMongo(nick string) ([]entities.Post, error) {
	filter := bson.M{"authorNick": nick}

	cursor, err := repository.collection.Find(context.TODO(), filter)
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
