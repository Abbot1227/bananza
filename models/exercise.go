package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Models

// TextExercise is type 0 is EnDe, 1 is DeEn
type TextExercise struct {
	ID       primitive.ObjectID `bson:"_id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
	Answer   string             `json:"answer"`
	Level    int                `json:"level"`
}

// ImageExercise is type 2
type ImageExercise struct {
	ID       primitive.ObjectID `bson:"_id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
	Answer   string             `json:"answer"`
	Level    int                `json:"level"`
}

// ImagesExercise is type 3
type ImagesExercise struct {
	ID     primitive.ObjectID    `bson:"_id"`
	Type   int                   `json:"type"`
	Word   string                `json:"word"`
	Cards  []ImagesExerciseCards `json:"cards"`
	Answer int                   `json:"answer"`
	Level  int                   `json:"level"`
}

// ImagesExerciseCards stores card for type 3
type ImagesExerciseCards struct {
	Word  string `json:"word"`
	Media string `json:"media"`
}

// AudioExercise is type 4
type AudioExercise struct {
	ID     primitive.ObjectID `bson:"_id"`
	Type   int                `json:"type"`
	Answer string             `json:"answer"`
	Level  int                `json:"level"`
}

// Forms

// 	Acquire

type AcquireExercise struct {
	Exp  int    `json:"exp"`
	Lang string `json:"lang"`
}

//	Send

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

//	Input

type InputAnswer struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
	User   string `json:"user"`
}
