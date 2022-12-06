package service

import (
	"Bananza/models"
)

type Authorization interface {
	AuthenticateUser(token models.AuthToken) (*models.User, error)
}

type Service struct {
	Authorization
}

//func NewService() *Service {
//	return &Service{Authorization: NewAuthService()}
//}
