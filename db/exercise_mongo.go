package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExerciseMongo struct {
}

func NewExerciseMongo() *ExerciseMongo {
	return &ExerciseMongo{}
}

func (r *ExerciseMongo) GetEnLnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 0}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", exerciseDesc.Exp}}}}}}
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

	if err = cursor.All(ctx, exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetLnEnExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.TextExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 1}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", exerciseDesc.Exp}}}}}}
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

	if err = cursor.All(ctx, exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetImageExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImageExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 2}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", exerciseDesc.Exp}}}}}}
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

	if err = cursor.All(ctx, exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetImagesExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.ImagesExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 3}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", exerciseDesc.Exp}}}}}}
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

	if err = cursor.All(ctx, exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetAudioExercise(ctx context.Context, exerciseDesc models.AcquireExercise, exercise *[]models.AudioExercise) error {
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 4}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", exerciseDesc.Exp}}}}}}
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

	if err = cursor.All(ctx, exercise); err != nil {
		return err
	}
	return nil
}

func (r *ExerciseMongo) GetRightAnswer(ctx context.Context, questionId string) (interface{}, error) {
	var answerStruct bson.D
	id, _ := primitive.ObjectIDFromHex(questionId)
	filter := bson.D{{"_id", id}}
	opts := options.FindOne().SetProjection(bson.D{{"_id", 0}, {"answer", 1}})

	if err := deExercisesCollection.FindOne(ctx, filter, opts).Decode(&answerStruct); err != nil {
		if err == mongo.ErrNoDocuments {
			if err = krExercisesCollection.FindOne(ctx, filter, opts).Decode(&answerStruct); err != nil {
				return "", err
			}
			answer := answerStruct.Map()

			return answer["answer"], nil
		}
		return "", err
	}
	answer := answerStruct.Map()

	return answer["answer"], nil
}

func (r *ExerciseMongo) IncrementProgressLevel(ctx context.Context, languageId string, expToAdd float64) error {
	id, _ := primitive.ObjectIDFromHex(languageId)

	result, err := userProgressCollection.UpdateByID(ctx, id, bson.D{
		{"$inc", bson.D{{"level", expToAdd}}},
	})
	if err != nil {
		return err
	}
	logrus.Println(result)

	result, err = userProgressCollection.UpdateByID(ctx, id, bson.D{
		{"$inc", bson.D{{"balance", 10}}},
	})
	if err != nil {
		return err
	}
	logrus.Println(result)

	return nil
}

func (r *ExerciseMongo) CreateTextImageExercise(ctx context.Context, exercise models.TextExercise, language string) error {
	objectId := primitive.NewObjectID()
	exercise.ID = objectId

	var languageCollection *mongo.Collection

	if language == "de" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	result, err := languageCollection.InsertOne(ctx, &exercise)
	if err != nil {
		return err
	}
	logrus.Println(result)

	return nil
}

func (r *ExerciseMongo) CreateImagesExercise(ctx context.Context, exercise models.ImagesExercise, language string) error {
	objectId := primitive.NewObjectID()
	exercise.ID = objectId

	var languageCollection *mongo.Collection

	if language == "de" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	result, err := languageCollection.InsertOne(ctx, &exercise)
	if err != nil {
		return err
	}
	logrus.Println(result)

	return nil
}

func (r *ExerciseMongo) CreateAudioExercise(ctx context.Context, exercise models.AudioExercise, language string) error {
	objectId := primitive.NewObjectID()
	exercise.ID = objectId

	var languageCollection *mongo.Collection

	if language == "de" {
		languageCollection = deExercisesCollection
	} else {
		languageCollection = krExercisesCollection
	}

	result, err := languageCollection.InsertOne(ctx, &exercise)
	if err != nil {
		return err
	}
	logrus.Println(result)

	return nil
}
