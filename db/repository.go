package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var usersCollection = OpenCollection(Client, "users")
var userProgressCollection = OpenCollection(Client, "userProgress")
var deExercisesCollection = OpenCollection(Client, "deExercises")
var krExercisesCollection = OpenCollection(Client, "krExercises")

type Authorization interface {
	FindUser(ctx context.Context, ID string, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) error
	FindLanguage(ctx context.Context, user *models.User, lastLanguageProgress *models.UserProgress) error
}

type User interface {
	AddLanguage(ctx context.Context, userProgress *models.UserProgress) (*mongo.InsertOneResult, error)
	FindUsers(ctx context.Context) (*[]models.User, error)
	FindUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error)
	FindProgress(ctx context.Context, userId primitive.ObjectID, language string) (*models.UserProgress, error)
	FindProgresses(ctx context.Context, userId primitive.ObjectID) (*[]models.UserProgress, error)
	SetProgressLevel(ctx context.Context, userProgressUpdate models.UserProgressUpdate) error
	SetLastLanguage(ctx context.Context, userId primitive.ObjectID, language string) error
}

type Exercise interface {
	GetEnLnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error
	GetLnEnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error
	GetImageExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImageExercise) error
	GetImagesExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImagesExercise) error
	GetAudioExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.AudioExercise) error
}

type Forum interface {
}

type Repository struct {
	Authorization
	User
	Exercise
	Forum
}

func NewRepository() *Repository {
	return &Repository{}
}
