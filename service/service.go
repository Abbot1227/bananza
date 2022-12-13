package service

import (
	"Bananza/db"
	"Bananza/models"
	"mime/multipart"

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
	UpdateProgress(languageId string, expToAdd float64) error
	CreateTextImageExercise(exercise models.TextExercise, language string) error
	CreateImagesExercise(exercise models.ImagesExercise, language string) error
	CreateAudioExercise(exercise models.AudioExercise, language string) error
	GetAudioAnswer(file multipart.File, language string) (interface{}, error)
}

type Forum interface {
	AddPost(inputForumPost *models.InputForumPost, forumPost *models.ForumPost) error
	GetForumTitles(skip int) (*[]models.SendForumTitles, error)
	GetForumPost(postId primitive.ObjectID) (*models.ForumPost, error)
	AddComment(inputComment *models.InputForumComment, postComment *models.ForumComment, postId primitive.ObjectID) error
	RemovePost(postId primitive.ObjectID) error
}

type Shop interface {
	BuyAvatar(inputAvatarPurchase *models.InputAvatarPurchase) error
	GetAvatars() (*[]models.Avatar, error)
	SetAvatar(inputAvatarSet *models.InputAvatarSet) error
}

type Grammar interface {
	GetGrammar(inputGrammar models.InputDictionary) (*[]models.Grammar, error)
	GetDictionary(inputDictionary models.InputDictionary) (*[]models.Dictionary, error)
}

type Service struct {
	Authorization
	User
	Exercise
	Forum
	Shop
	Grammar
}

func NewService(repos *db.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
		Exercise:      NewExerciseService(repos.Exercise),
		Forum:         NewForumService(repos.Forum),
		Shop:          NewShopService(repos.Shop),
		Grammar:       NewGrammarService(repos.Grammar),
	}
}
