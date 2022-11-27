package models

type TextExercise struct {
	ID       string
	Question string
	Answer   string
	Level    int
}

type AudioExercise struct {
	ID    string
	Text  string
	Level int
}

type ExercisesSet struct {
	TextExercises map[string]TextExercise
	AudioExercise map[string]AudioExercise
	Level         int
}

type InputTextExercise struct {
}
