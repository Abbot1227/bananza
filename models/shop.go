package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Models

type Avatar struct {
	Id    primitive.ObjectID `json:"-" bson:"_id"`
	Url   string             `json:"url"`
	Price int                `json:"price"`
}

// Forms

type InputAvatarPurchase struct {
	UserId    string `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
	Price     int    `json:"price"`
}

type InputAvatarSet struct {
	UserId    string `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
}
