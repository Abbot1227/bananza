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
		return nil, err
	}
	defer cancel()

	return result, nil
}

func (s *UserService) FindProfiles() (*[]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	users, err := s.repo.FindUsers(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return users, nil
}

func (s *UserService) FindProfile(userId primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	user, err := s.repo.FindUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return user, nil
}

func (s *UserService) FindProgress(userId primitive.ObjectID, language string) (*models.UserProgress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	progress, err := s.repo.FindProgress(ctx, userId, language)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return progress, nil
}

func (s *UserService) FindProgresses(userId primitive.ObjectID) (*[]models.UserProgress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	progresses, err := s.repo.FindProgresses(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return progresses, nil
}

func (s *UserService) SetProgressLevel(userProgressUpdate models.UserProgressUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.SetProgressLevel(ctx, userProgressUpdate); err != nil {
		return err
	}
	defer cancel()

	return nil
}

func (s *UserService) SetLastLanguage(userId primitive.ObjectID, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.SetLastLanguage(ctx, userId, language); err != nil {
		return err
	}
	defer cancel()

	return nil
}
