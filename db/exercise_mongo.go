package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseMongo struct {
}

func NewExerciseMongo() *ExerciseMongo {
	return &ExerciseMongo{}
}

func (r *ExerciseMongo) GetEnLnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 1}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", exerciseDesc.Exp}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	var languageCollection *mongo.Collection

	if exerciseDesc.Lang == "German" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	cursor, err := languageCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, &exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetLnEnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 2}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", exerciseDesc.Exp}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	var languageCollection *mongo.Collection

	if exerciseDesc.Lang == "German" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	cursor, err := languageCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, &exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetImageExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImageExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 2}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", exerciseDesc.Exp}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	var languageCollection *mongo.Collection

	if exerciseDesc.Lang == "German" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	cursor, err := languageCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, &exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetImagesExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImagesExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 3}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", exerciseDesc.Exp}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	var languageCollection *mongo.Collection

	if exerciseDesc.Lang == "German" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	cursor, err := languageCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, &exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetAudioExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.AudioExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 4}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", exerciseDesc.Exp}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	var languageCollection *mongo.Collection

	if exerciseDesc.Lang == "German" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	cursor, err := languageCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, &exercise); err != nil {
		return err
	}
	return nil
}
