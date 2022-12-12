package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopMongo struct {
}

func NewShopMongo() *ShopMongo {
	return &ShopMongo{}
}

func (r *ShopMongo) GetUser(ctx context.Context, userId string) (models.User, error) {
	var user models.User
	id, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.D{{"_id", id}}

	if err := usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}
	logrus.Println(user)
	return user, nil
}

func (r *ShopMongo) UpdateUserBalance(ctx context.Context, user *models.User) error {
	update := bson.D{{"$set", bson.D{{"balance", user.Balance}}}}

	result, err := usersCollection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}

func (r *ShopMongo) AddAvatarToUser(ctx context.Context, userId string, avatarUrl string) error {
	id, _ := primitive.ObjectIDFromHex(userId)
	update := bson.D{{"$push", bson.D{{"avatars", avatarUrl}}}}

	result, err := usersCollection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}

func (r *ShopMongo) SetUserAvatar(ctx context.Context, userId string, avatarUrl string) error {
	id, _ := primitive.ObjectIDFromHex(userId)
	update := bson.D{{"$set", bson.D{{"avatarurl", avatarUrl}}}}

	result, err := usersCollection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}

func (r *ShopMongo) GetAvatars(ctx context.Context) ([]models.Avatar, error) {
	var avatars []models.Avatar

	cursor, err := avatarsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &avatars); err != nil {
		return nil, err
	}
	logrus.Println(avatars)
	return avatars, nil
}
