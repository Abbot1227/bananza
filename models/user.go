package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	UserId    string             `json:"user_id"`
	AvatarURL string             `json:"avatar_url"`
}

type UserProgress struct {
	ID       primitive.ObjectID `bson:"_id"`
	Language string             `json:"language"`
	Level    int                `json:"level"`
	User     primitive.ObjectID `bson:"user"`
}
