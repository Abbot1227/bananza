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
var postsCollection = OpenCollection(Client, "posts")

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
	DeleteProfile(ctx context.Context, userId primitive.ObjectID) error
}

type Exercise interface {
	GetEnLnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error
	GetLnEnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error
	GetImageExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImageExercise) error
	GetImagesExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImagesExercise) error
	GetAudioExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.AudioExercise) error
	GetRightAnswer(ctx context.Context, questionId string) (interface{}, error)
	IncrementProgressLevel(ctx context.Context, languageId string, expToAdd float64) error
	CreateTextImageExercise(ctx context.Context, exercise models.TextExercise, language string) error
	CreateImagesExercise(ctx context.Context, exercise models.ImagesExercise, language string) error
	CreateAudioExercise(ctx context.Context, exercise models.AudioExercise, language string) error
}

type Forum interface {
	FindUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error)
	CreatePost(ctx context.Context, forumPost *models.ForumPost) error
	GetForumPosts(ctx context.Context) ([]models.ForumPost, error)
	GetForumPost(ctx context.Context, postId primitive.ObjectID) (*models.ForumPost, error)
	CreateComment(ctx context.Context, forumComment *models.ForumComment, postId primitive.ObjectID) error
	DeletePost(ctx context.Context, postId primitive.ObjectID) error
}

type Repository struct {
	Authorization
	User
	Exercise
	Forum
}

func NewRepository() *Repository {
	return &Repository{
		Authorization: NewAuthMongo(),
		User:          NewUserMongo(),
		Exercise:      NewExerciseMongo(),
		Forum:         NewForumMongo(),
	}
}
