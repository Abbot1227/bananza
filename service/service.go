package service

import (
	"Bananza/db"
	"Bananza/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	AuthenticateUser(token models.AuthToken) (*models.User, error)
	GetLastLanguage(user *models.User) (*models.UserProgress, error)
}

type User interface {
	AddLanguage(inputLanguage models.InputLanguage) (*mongo.InsertOneResult, error)
	FindProfile(userId primitive.ObjectID) (*models.User, error)
}

type Exercise interface {
}

type Forum interface {
}

type Service struct {
	Authorization
	User
	Exercise
	Forum
}

func NewService() *Service {
	return &Service{Authorization: NewAuthService(db.NewAuthMongo()),
		User: NewUserService(db.NewUserMongo())}
}
