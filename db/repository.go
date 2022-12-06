package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection = OpenCollection(Client, "users")

type Authorization interface {
	FindUser(ctx context.Context, ID string, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
}

type Repository struct {
	Authorization
}

func NewRepository() *Repository {
	return &Repository{}
}
