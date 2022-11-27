package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
)

type AuthService struct {
	repo db.Authorization
}

func NewAuthService(repo db.Authorization) *AuthService {
	return &AuthService{}
}

func (s *AuthService) ValidateToken(ctx context.Context, token models.AuthToken) (*models.User, error) {
	return nil, nil
}

func (s *AuthService) AuthenticateUser(user models.User) error {
	return nil
}
