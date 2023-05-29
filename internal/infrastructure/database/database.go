package database

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Database, error) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(os.Getenv("DB_MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(os.Getenv("DB_NAME")), nil
}
