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

func (r *UserMongo) FindUsers(ctx context.Context) (*[]models.User, error) {
	var users []models.User

	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	logrus.Println(users)
	return &users, nil
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

func (r *UserMongo) FindProgress(ctx context.Context, userId primitive.ObjectID, language string) (*models.UserProgress, error) {
	var languageProgress models.UserProgress
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user", userId}},
				bson.D{{"language", language}},
			},
		},
	}

	if err := userProgressCollection.FindOne(ctx, filter).Decode(&languageProgress); err != nil {
		return nil, err
	}
	logrus.Println(languageProgress)
	return &languageProgress, nil
}

func (r *UserMongo) FindProgresses(ctx context.Context, userId primitive.ObjectID) (*[]models.UserProgress, error) {
	var languagesProgress []models.UserProgress
	filter := bson.D{{"user", userId}}

	cursor, err := usersCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &languagesProgress); err != nil {
		return nil, err
	}
	logrus.Println(languagesProgress)
	return &languagesProgress, nil
}

func (r *UserMongo) SetProgressLevel(ctx context.Context, userProgressUpdate models.UserProgressUpdate) error {
	update := bson.D{
		{"$set",
			bson.D{
				{"level", userProgressUpdate.Level},
			}},
	}

	result, err := userProgressCollection.UpdateByID(ctx, userProgressUpdate.ProgressId, update)
	if err != nil {
		return nil
	}
	logrus.Println(result)
	return nil
}

func (r *UserMongo) SetLastLanguage(ctx context.Context, userId primitive.ObjectID, language string) error {
	update := bson.D{
		{"$set",
			bson.D{
				{"lastlanguage", language},
			},
		},
	}

	result, err := usersCollection.UpdateByID(ctx, userId, update)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}
