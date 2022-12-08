package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type ExerciseService struct {
	repo   db.Exercise
	client *http.Client
}

func NewExerciseService(repo db.Exercise) *ExerciseService {
	return &ExerciseService{repo: repo,
		client: new(http.Client)}
}

func (s *ExerciseService) GetExerciseType() (int, error) {
	if err := s.checkASRConnection(); err != nil {
		return 0, nil
	}

	exerciseType := s.generateRandomType(4)
	return exerciseType, nil
}

func (s *ExerciseService) GetEnLnExercise(exerciseDesc models.AcquireExercise, exercise *models.SendTextExercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var exercises []models.TextExercise

	if err := s.repo.GetEnLnExercise(ctx, exerciseDesc, &exercises); err != nil {
		return err
	}
	defer cancel()

	// Assigning corresponding fields of random exercise to the exercise to be sent
	for _, randExercise := range exercises {
		exercise.ID = randExercise.ID
		exercise.Type = randExercise.Type
		exercise.Question = randExercise.Question
	}
	return nil
}

func (s *ExerciseService) GetLnEnExercise(exerciseDesc models.AcquireExercise, exercise *models.SendTextExercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var exercises []models.TextExercise

	if err := s.repo.GetLnEnExercise(ctx, exerciseDesc, &exercises); err != nil {
		return err
	}
	defer cancel()

	// Assigning corresponding fields of random exercise to the exercise to be sent
	for _, randExercise := range exercises {
		exercise.ID = randExercise.ID
		exercise.Type = randExercise.Type
		exercise.Question = randExercise.Question
	}
	return nil
}

func (s *ExerciseService) GetImageExercise(exerciseDesc models.AcquireExercise, exercise *models.SendImageExercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var exercises []models.ImageExercise

	if err := s.repo.GetImageExercise(ctx, exerciseDesc, &exercises); err != nil {
		return err
	}
	defer cancel()

	// Assigning corresponding fields of random exercise to the exercise to be sent
	for _, randExercise := range exercises {
		exercise.ID = randExercise.ID
		exercise.Type = randExercise.Type
		exercise.Question = randExercise.Question
	}
	return nil
}

func (s *ExerciseService) GetImagesExercise(exerciseDesc models.AcquireExercise, exercise *models.SendImagesExercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var exercises []models.ImagesExercise

	if err := s.repo.GetImagesExercise(ctx, exerciseDesc, &exercises); err != nil {
		return err
	}
	defer cancel()

	// Assigning corresponding fields of random exercise to the exercise to be sent
	for _, randExercise := range exercises {
		exercise.ID = randExercise.ID
		exercise.Type = randExercise.Type
		exercise.Word = randExercise.Word
		exercise.Cards = randExercise.Cards
	}
	return nil
}

func (s *ExerciseService) GetAudioExercise(exerciseDesc models.AcquireExercise, exercise *models.SendAudioExercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var exercises []models.AudioExercise

	if err := s.repo.GetAudioExercise(ctx, exerciseDesc, &exercises); err != nil {
		return err
	}
	defer cancel()

	// Assigning corresponding fields of random exercise to the exercise to be sent
	for _, randExercise := range exercises {
		exercise.ID = randExercise.ID
		exercise.Type = randExercise.Type
		exercise.Question = randExercise.Answer
	}
	return nil
}

func (s *ExerciseService) GetRightAnswer(questionId string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	answer, err := s.repo.GetRightAnswer(ctx, questionId)
	if err != nil {
		return "", err
	}
	defer cancel()

	return answer, nil
}

func (s *ExerciseService) UpdateProgress(languageId string, expToAdd int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.IncrementProgressLevel(ctx, languageId, expToAdd); err != nil {
		return err
	}
	defer cancel()

	return nil
}

func (s *ExerciseService) CreateTextImageExercise(exercise models.TextExercise, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.CreateTextImageExercise(ctx, exercise, language); err != nil {
		return err
	}
	defer cancel()

	return nil
}

func (s *ExerciseService) CreateImagesExercise(exercise models.ImagesExercise, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.CreateImagesExercise(ctx, exercise, language); err != nil {
		return err
	}
	defer cancel()

	return nil
}

func (s *ExerciseService) CreateAudioExercise(exercise models.AudioExercise, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.CreateAudioExercise(ctx, exercise, language); err != nil {
		return err
	}
	defer cancel()

	return nil
}

func (s *ExerciseService) checkASRConnection() error {
	req, err := http.NewRequest("GET", "http://localhost:4040/predict", nil)
	if err != nil {
		return err
	}

	res, _ := s.client.Do(req)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}
	return nil
}

// generateRandomType is a function that generates number between 0 and 4
// which is a type of exercise
func (s *ExerciseService) generateRandomType(last int) int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := last

	return rand.Intn(max-min+1) + min
}
