package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongo struct {
}

func NewUserMongo() *UserMongo {
	return &UserMongo{}
}

func (r *UserMongo) AddLanguage(ctx context.Context, userProgress *models.UserProgress) (*mongo.InsertOneResult, error) {
	result, err := userProgressCollection.InsertOne(ctx, userProgress)
	if err != nil {
		return nil, err
	}
	logrus.Println(result)
	return result, nil
}
