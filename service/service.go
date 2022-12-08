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
	FindProfiles() (*[]models.User, error)
	FindProfile(userId primitive.ObjectID) (*models.User, error)
	FindProgress(userId primitive.ObjectID, language string) (*models.UserProgress, error)
	FindProgresses(userId primitive.ObjectID) (*[]models.UserProgress, error)
	SetProgressLevel(userProgressUpdate models.UserProgressUpdate) error
	SetLastLanguage(userId primitive.ObjectID, language string) error
	DeleteProfile(userId primitive.ObjectID) error
}

type Exercise interface {
	GetExerciseType() (int, error)
	GetEnLnExercise(exerciseDesc models.AcquireExercise, exercise *models.SendTextExercise) error
	GetLnEnExercise(exerciseDesc models.AcquireExercise, exercise *models.SendTextExercise) error
	GetImageExercise(exerciseDesc models.AcquireExercise, exercise *models.SendImageExercise) error
	GetImagesExercise(exerciseDesc models.AcquireExercise, exercise *models.SendImagesExercise) error
	GetAudioExercise(exerciseDesc models.AcquireExercise, exercise *models.SendAudioExercise) error
	GetRightAnswer(questionId string) (interface{}, error)
	UpdateProgress(languageId string, expToAdd int) error
	CreateTextImageExercise(exercise models.TextExercise, language string) error
	CreateImagesExercise(exercise models.ImagesExercise, language string) error
	CreateAudioExercise(exercise models.AudioExercise, language string) error
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
