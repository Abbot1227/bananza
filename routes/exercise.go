package routes

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"net/http"
	"time"
)

var tempExercisesCollection = db.OpenCollection(db.Client, "tempExercises")

func LangExercise(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var acquireExercise models.AcquireExercise

	fmt.Println(acquireExercise)

	if err := c.BindJSON(&acquireExercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&acquireExercise)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	exerciseType := generateRandomType()
	level := acquireExercise.Exp

	if acquireExercise.Lang == "de" {
		// TODO do nothing
	}

	fmt.Println(exerciseType)

	switch exerciseType {
	case 0:
		var exercise models.SendTextExercise

		if err := sendEnDeExercise(ctx, level, &exercise); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()
			break
		}
		fmt.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	case 1:
		var exercise models.SendTextExercise

		if err := sendDeEnExercise(ctx, level, &exercise); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()
			break
		}
		fmt.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	case 2:
		var exercise models.SendImageExercise

		if err := sendImageExercise(ctx, level, &exercise); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()
			break
		}
		fmt.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	//case 3:
	//	if err := sendImagesExercise(ctx, level); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		defer cancel()
	//		break
	//	}
	case 4:
		var exercise models.SendAudioExercise

		if err := sendAudioExercise(ctx, level, &exercise); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()
			break
		}
		fmt.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		defer cancel()
		break
	}
	defer cancel()
}

func sendEnDeExercise(ctx context.Context, level int, sendExercise *models.SendTextExercise) error {
	var exercise models.TextExercise
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"type", 0}},
				bson.D{{"level", level}},
			},
		},
	}

	if err := tempExercisesCollection.FindOne(ctx, filter).Decode(&exercise); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Got Exercise:", exercise)

	sendExercise.ID = exercise.ID
	sendExercise.Type = exercise.Type
	sendExercise.Question = exercise.Question

	return nil
}

func sendDeEnExercise(ctx context.Context, level int, sendExercise *models.SendTextExercise) error {
	var exercise models.TextExercise
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"type", 1}},
				bson.D{{"level", level}},
			},
		},
	}

	if err := tempExercisesCollection.FindOne(ctx, filter).Decode(&exercise); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Got Exercise:", exercise)

	sendExercise.ID = exercise.ID
	sendExercise.Type = exercise.Type
	sendExercise.Question = exercise.Question

	return nil
}

func sendImageExercise(ctx context.Context, level int, sendExercise *models.SendImageExercise) error {
	var exercise models.ImageExercise
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"type", 2}},
				bson.D{{"level", level}},
			},
		},
	}

	if err := tempExercisesCollection.FindOne(ctx, filter).Decode(&exercise); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(exercise)

	sendExercise.ID = exercise.ID
	sendExercise.Type = exercise.Type
	sendExercise.Question = exercise.Question

	return nil
}

//func sendImagesExercise(ctx context.Context, level int) error {
//	var exercise models.ImagesExercise
//	filter := bson.D{
//		{"$and",
//			bson.A{
//				bson.D{{"type", 3}},
//				bson.D{{"level", level}},
//			},
//		},
//	}
//}

func sendAudioExercise(ctx context.Context, level int, sendExercise *models.SendAudioExercise) error {
	var exercise models.AudioExercise
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"type", 4}},
				bson.D{{"level", level}},
			},
		},
	}

	if err := tempExercisesCollection.FindOne(ctx, filter).Decode(&exercise); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(exercise)

	sendExercise.ID = exercise.ID
	sendExercise.Type = exercise.Type
	sendExercise.Question = exercise.Answer

	return nil
}

func SendAnswer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	language := c.Params.ByName("lang")
	language = language[5:]

	var inputAnswer models.InputAnswer

	if err := c.BindJSON(&inputAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&inputAnswer)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	var answerStruct bson.D
	questionId, _ := primitive.ObjectIDFromHex(inputAnswer.ID)
	filter := bson.D{{"_id", questionId}}
	opts := options.FindOne().SetProjection(bson.D{{"_id", 0}, {"answer", 1}})

	if err := tempExercisesCollection.FindOne(ctx, filter, opts).Decode(&answerStruct); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		defer cancel()
		return
	}
	defer cancel()
	answer := answerStruct.Map()

	fmt.Println(questionId)
	fmt.Println("Right:", answer["answer"])
	fmt.Println("User:", inputAnswer.Answer)

	// Добавить прибавление очков пользователю за правильный ответ
	if inputAnswer.Answer == answer["answer"] {
		c.JSON(http.StatusOK, gin.H{"correct": "true", "answer": answer["answer"]})
	} else {
		c.JSON(http.StatusOK, gin.H{"correct": "false", "answer": answer["answer"]})
		return
	}

	languageId, _ := primitive.ObjectIDFromHex(inputAnswer.LanguageId)

	result, err := userProgressCollection.UpdateByID(ctx, languageId, bson.D{
		{"$inc", bson.D{{"level", 15}}},
	})
	if err != nil {
		fmt.Println("Could not add points to user")
	}
	fmt.Println(result)
}

// generateRandomType is a function that generates number between 0 and 4
// which is a type of exercise
func generateRandomType() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 4

	return rand.Intn(max-min+1) + min
}

func calculateGainExp(level int) int {
	return 1/level - 1
}
