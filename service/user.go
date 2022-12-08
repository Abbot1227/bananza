package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserService struct {
	repo db.User
}

func NewUserService(repo db.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddLanguage(inputLanguage models.InputLanguage) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	userId, _ := primitive.ObjectIDFromHex(inputLanguage.User)

	userProgress := &models.UserProgress{ID: primitive.NewObjectID(),
		Language: inputLanguage.Language,
		Level:    0,
		User:     userId}

	result, err := s.repo.AddLanguage(ctx, userProgress)
	if err != nil {
		return result, err
	}
	defer cancel()

	return result, err
}
