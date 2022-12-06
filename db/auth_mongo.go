package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
}

func NewAuthMongo(db *mongo.Client) *AuthMongo {
	return &AuthMongo{}
}

func (r *AuthMongo) FindUser(ctx context.Context, ID string, user *models.User) error {
	// Filter to find user with specified Google ID in users collection
	filter := bson.D{{"userid", ID}}

	// Obtain user info from users collection and store it in user object
	// if not found return error otherwise return nil
	if err := usersCollection.FindOne(ctx, filter).Decode(user); err != nil {
		logrus.Error(err.Error())
		return err
	}

	return nil
}

func (r *AuthMongo) CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	return nil, nil
}
