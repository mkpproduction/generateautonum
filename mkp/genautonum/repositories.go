package genautonum

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB *sql.DB
	MongoDB *mongo.Database
	Context context.Context
}

func NewRepository(conn *sql.DB, ctx context.Context, MongoDB *mongo.Database) Repository {
	return Repository{
		DB: conn,
		Context: ctx,
		MongoDB: MongoDB,
	}
}