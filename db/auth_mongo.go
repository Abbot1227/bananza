package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthMongo struct {
}

func NewAuthMongo() *AuthMongo {
	return &AuthMongo{}
}

func (r *AuthMongo) FindUser(ctx context.Context, ID string, user *models.User) error {
	// Filter to find user with specified Google ID in users collection
	filter := bson.D{{"userid", ID}}

	// Obtain user info from users collection and store it in user object
	// if not found return error otherwise return nil
	if err := usersCollection.FindOne(ctx, filter).Decode(user); err != nil {
		return err
	}
	return nil
}

func (r *AuthMongo) CreateUser(ctx context.Context, user *models.User) error {
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}

func (r *AuthMongo) FindLanguage(ctx context.Context, user *models.User, lastLanguageProgress *models.UserProgress) error {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user", user.ID}},
				bson.D{{"language", user.LastLanguage}},
			},
		},
	}

	if err := userProgressCollection.FindOne(ctx, filter).Decode(lastLanguageProgress); err != nil {
		return err
	}
	return nil
}
