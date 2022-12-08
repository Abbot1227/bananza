package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *UserMongo) FindUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error) {
	var user models.User
	filter := bson.D{{"_id", userId}}

	if err := usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	logrus.Println(user)
	return &user, nil
}
