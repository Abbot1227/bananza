package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Models

type ForumPost struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `json:"title"`
	Text      string             `json:"text"`
	Author    string             `json:"author"`
	CreatedAt primitive.DateTime `json:"created_at"`
	Replies   []ForumComment     `json:"replies"`
}

type ForumComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `json:"text"`
	Author    string             `json:"author"`
	CreatedAt primitive.DateTime `json:"created_at"`
}

// Forms

type SendForumTitles struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `json:"title"`
	CreatedAt primitive.DateTime `json:"created_at"`
}

type InputForumPost struct {
	UserId string `json:"user_id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type InputForumComment struct {
	UserId string `json:"user_id"`
	Text   string `json:"text"`
}
