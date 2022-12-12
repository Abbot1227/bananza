package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"errors"
	"time"
)

type ShopService struct {
	repo db.Shop
}

func NewShopService(repo db.Shop) *ShopService {
	return &ShopService{repo: repo}
}

func (s *ShopService) BuyAvatar(inputAvatarPurchase *models.InputAvatarPurchase) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Fetching user document from db to subtract price of avatar from user's balance
	user, err := s.repo.GetUser(ctx, inputAvatarPurchase.UserId)
	if err != nil {
		return err
	}
	defer cancel()

	for _, avatar := range user.Avatars {
		if inputAvatarPurchase.AvatarUrl == avatar {
			return errors.New("user already has this avatar")
		}
	}

	// Checking if user has sufficient amount of money to pay for avatar
	if user.Balance < inputAvatarPurchase.Price {
		return errors.New("insufficient balance")
	} else {
		user.Balance -= inputAvatarPurchase.Price
	}

	// Subtracting price of avatar from user's balance
	if err = s.repo.UpdateUserBalance(ctx, &user); err != nil {
		return err
	}

	// Adding avatar to array of user's avatars, so later he can choose between them
	if err = s.repo.AddAvatarToUser(ctx, inputAvatarPurchase.UserId, inputAvatarPurchase.AvatarUrl); err != nil {
		return err
	}

	// Changing user's current avatar
	if err = s.repo.SetUserAvatar(ctx, inputAvatarPurchase.UserId, inputAvatarPurchase.AvatarUrl); err != nil {
		return err
	}

	return nil
}

func (s *ShopService) GetAvatars() (*[]models.Avatar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var avatars []models.Avatar

	avatars, err := s.repo.GetAvatars(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return &avatars, nil
}

func (s *ShopService) SetAvatar(inputAvatarSet *models.InputAvatarSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Fetching user document from db to check whether user has avatar or not
	user, err := s.repo.GetUser(ctx, inputAvatarSet.UserId)
	if err != nil {
		return err
	}

	for _, avatar := range user.Avatars {
		if inputAvatarSet.AvatarUrl == avatar {
			if err := s.repo.SetUserAvatar(ctx, inputAvatarSet.UserId, inputAvatarSet.AvatarUrl); err != nil {
				return err
			}
			return nil
		}
	}
	defer cancel()

	return errors.New("user does not have this avatar")
}
