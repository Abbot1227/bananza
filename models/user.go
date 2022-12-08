package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Models

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Email        string             `json:"email"`
	Name         string             `json:"name"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	UserId       string             `json:"user_id"`
	AvatarURL    string             `json:"avatar_url"`
	LastLanguage string             `json:"last_language"`
	Balance      int                `json:"balance"`
}

type UserProgress struct {
	ID       primitive.ObjectID `bson:"_id"`
	Language string             `json:"language"`
	User     primitive.ObjectID `bson:"user"`
	Level    int                `json:"level"`
}

// Forms

type InputLanguage struct {
	Language string `json:"language"`
	User     string `json:"user"`
}

type LanguageUpdate struct {
	Language string `json:"language"`
	Level    int    `json:"level"`
	User     string `json:"user"`
}

type UserProgressUpdate struct {
	ProgressId string `json:"progress_id"`
	Level      int    `json:"level"`
}
