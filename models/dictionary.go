package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Models

type Dictionary struct {
	ID       primitive.ObjectID `json:"-" bson:"_id"`
	A        string             `json:"a" bson:"a"` // German/Korean
	B        string             `json:"b" bson:"b"` // English
	Level    int                `json:"level" bson:"level"`
	Language string             `json:"language" bson:"language"` // "German" or "Korean"
}

type Grammar struct {
	ID       primitive.ObjectID `json:"-" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Text     string             `json:"text" bson:"text"`
	Level    int                `json:"level" bson:"level"`
	Language string             `json:"language" bson:"language"` // "German" or "Korean"
}

// Forms

type InputDictionary struct {
	Level    int    `json:"level" bson:"level"`
	Language string `json:"language" bson:"language"` // "German" or "Korean"
}
