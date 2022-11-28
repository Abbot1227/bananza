package routes

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"net/http"
	"time"
)

var tempExercisesCollection = db.OpenCollection(db.Client, "tempExercises")

func LangExercise(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var acquireExercise models.AcquireExercise

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
	//case 4:
	//	sendAudioExercise(ctx, level)
	//	if err := sendAudioExercise(ctx, level); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//		defer cancel()
	//		break
	//	}
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
//
//func sendAudioExercise(ctx context.Context, level int) error {
//	var exercise models.AudioExercise
//	filter := bson.D{
//		{"$and",
//			bson.A{
//				bson.D{{"type", 4}},
//				bson.D{{"level", level}},
//			},
//		},
//	}
//}

func SendAnswer(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	//
	//var inputAnswer models.InputAnswer
	//
	//if err := c.BindJSON(&inputAnswer); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	fmt.Println(err)
	//	defer cancel()
	//	return
	//}
	//
	//// Ensure that data we receive is correct
	//validationErr := validate.Struct(&inputAnswer)
	//if validationErr != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
	//	defer cancel()
	//	return
	//}
	//
	//if err := tempExercisesCollection.FindOne(ctx)

}

// generateRandomType is a function that generates number between 0 and 4
// which is a type of exercise
func generateRandomType() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 2

	return rand.Intn(max-min+1) + min
}
