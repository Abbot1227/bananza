package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TextExercise struct {
	ID       primitive.ObjectID `bson:"id"`
	Type     int                `json:"type"`
	Question string             `json:"question"`
	Answer   string             `json:"answer"`
	Level    int                `json:"level"`
}

type AudioExercise struct {
	ID    string
	Text  string
	Level int
}

type ExercisesSet struct {
	TextExercises map[int]TextExercise
	AudioExercise map[int]AudioExercise
	Level         int
}

type InputTextExercise struct {
	ID     primitive.ObjectID `bson:"id"`
	Answer string             `json:"answer"`
	User   string             `json:"user"`
}
