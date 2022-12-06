package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
	"time"
)

type AuthService struct {
	repo db.Authorization
}

func NewAuthService(repo db.Authorization) *AuthService {
	return &AuthService{}
}

func (s *AuthService) AuthenticateUser(token models.AuthToken) (*models.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var user *models.User

	// Getting user google account info from Google api
	userInfo, err := getUserInfo(token.Token)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	userId := userInfo.Id

	if err = s.repo.FindUser(ctx, userId, user); err != nil {
		if err != mongo.ErrNoDocuments {
			logrus.Error(err.Error())
			return nil, err
		}

		// Getting user google account info from Google api
		userInfo, err := getUserInfo(token.Token)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}

		// Assigning userInfo fields from google API to newly created user
		user.ID = primitive.NewObjectID()
		user.Email = userInfo.Email
		user.Name = userInfo.Name
		user.FirstName = userInfo.GivenName
		user.LastName = userInfo.FamilyName
		user.UserId = userInfo.Id
		user.AvatarURL = userInfo.Picture
		user.LastLanguage = ""
		user.Balance = 0

		// Inserting new user into database

	}
	return nil, nil
}

// getUserInfo is a function TODO add description
func getUserInfo(token string) (*oauth2.Userinfoplus, error) {
	oauth2Service, err := oauth2.NewService(context.Background(), option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}
	userInfoService := oauth2.NewUserinfoV2MeService(oauth2Service)

	// Getting user google account info from Google api
	userInfo, err := userInfoService.Get().Do(googleapi.QueryParameter("access_token", token))
	if err != nil {
		e, _ := err.(*googleapi.Error)
		return nil, e
	}
	return userInfo, nil
}
