package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db *mongo.Client
}

func NewAuthMongo(db *mongo.Client) *AuthMongo {
	return &AuthMongo{db: db}
}

func (r *AuthMongo) FindUser(ctx context.Context, ID string) error {
	return nil
}

func (r *AuthMongo) CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	return nil, nil
}
