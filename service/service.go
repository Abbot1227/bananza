package service

import (
	"Bananza/models"
	"context"
)

type Authorization interface {
	ValidateToken(ctx context.Context, token models.AuthToken) (*models.User, error)
	AuthenticateUser(user models.User) error
}

type Service struct {
	Authorization
}

//func NewService() *Service {
//	return &Service{Authorization: NewAuthService()}
//}
