package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	FindUser(ctx context.Context, ID string) error
	CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
}

type Repository struct {
	Authorization
}

func NewRepository() *Repository {
	return &Repository{}
}
