package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Models

// TextExercise is type 0 is EnDe, 1 is DeEn
type TextExercise struct {
	ID       primitive.ObjectID `bson:"id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
	Answer   string             `json:"answer"`
	Level    int                `json:"level"`
}

// ImageExercise is type 2
type ImageExercise struct {
	ID       primitive.ObjectID `bson:"id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
	Answer   string             `json:"answer"`
	Level    int                `json:"level"`
}

// ImagesExercise is type 3
type ImagesExercise struct {
	ID     primitive.ObjectID    `bson:"id"`
	Type   int                   `json:"type"`
	Word   string                `json:"word"`
	Cards  []ImagesExerciseCards `json:"cards"`
	Answer int                   `json:"answer"`
	Level  int                   `json:"level"`
}

type ImagesExerciseCards struct {
	Word  string `json:"word"`
	Media string `json:"media"`
}

type AudioExercise struct {
	ID    primitive.ObjectID `bson:"id"`
	Type  int                `json:"type"`
	Word  string             `json:"word"`
	Level int                `json:"level"`
}

// Forms

type AcquireExercise struct {
	Exp  int    `json:"exp"`
	Lang string `json:"lang"`
}

type SendTextExercise struct {
	ID       primitive.ObjectID `bson:"id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
}

type SendImageExercise struct {
	ID       primitive.ObjectID `bson:"id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
}

type InputTextExercise struct {
	ID     primitive.ObjectID `bson:"id"`
	Answer string             `json:"answer"`
	User   string             `json:"user"`
	Level  int                `json:"level"`
}

type InputImageExercise struct {
	ID     primitive.ObjectID `bson:"id"`
	Answer string             `json:"answer"`
	User   string             `json:"user"`
	Level  int                `json:"level"`
}
