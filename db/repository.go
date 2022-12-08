package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection = OpenCollection(Client, "users")
var userProgressCollection = OpenCollection(Client, "userProgress")

type Authorization interface {
	FindUser(ctx context.Context, ID string, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) error
	FindLanguage(ctx context.Context, user *models.User, lastLanguageProgress *models.UserProgress) error
}

type User interface {
	AddLanguage(ctx context.Context, userProgress *models.UserProgress) (*mongo.InsertOneResult, error)
	FindUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error)
}

type Exercise interface {
}

type Repository struct {
	Authorization
	User
	Exercise
}

func NewRepository() *Repository {
	return &Repository{}
}
