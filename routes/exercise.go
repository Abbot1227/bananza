package routes

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"net/http"
	"time"
)

var ExpMultiplier = 15
var tempExercisesCollection = db.OpenCollection(db.Client, "tempExercises")
var deExercisesCollection = db.OpenCollection(db.Client, "deExercises")
var krExercisesCollection = db.OpenCollection(db.Client, "krExercises")

func SendExercise(c *gin.Context) {
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
	level := acquireExercise.Exp / 100

	// Disabled during test period since there is no connection with ASR model API
	//if err := checkASRConnection(); err != nil {
	//	exerciseType = 0
	//}

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
	case 3:
		var exercise models.SendImagesExercise

		if err := sendImagesExercise(ctx, level, &exercise); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()
			break
		}
		fmt.Println(exercise)

		c.JSON(http.StatusOK, exercise)
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

func SendAnswer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

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

	expToAdd := calculateGainExp(inputAnswer.Level)

	if inputAnswer.Answer == answer["answer"] {
		c.JSON(http.StatusOK, gin.H{"correct": "true", "answer": answer["answer"], "exp": expToAdd})
	} else {
		c.JSON(http.StatusOK, gin.H{"correct": "false", "answer": answer["answer"], "exp": 0})
		return
	}

	languageId, _ := primitive.ObjectIDFromHex(inputAnswer.LanguageId)

	result, err := userProgressCollection.UpdateByID(ctx, languageId, bson.D{
		{"$inc", bson.D{{"level", expToAdd}}},
	})
	if err != nil {
		fmt.Println("Could not add points to user")
	}
	fmt.Println(result)
}

// AddTextImageExercise is a function todo
func AddTextImageExercise(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	var exercise models.TextExercise

	if err := c.BindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&exercise)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	objectId := primitive.NewObjectID()
	exercise.ID = objectId

	switch language {
	case "de":
		defer cancel()
		result, err := deExercisesCollection.InsertOne(ctx, &exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exercise was not added"})
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	case "kr":
		defer cancel()
		result, err := krExercisesCollection.InsertOne(ctx, &exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exercise was not added"})
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	default: // If wrong exercise was provided
		fmt.Println(exercise)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect exercise"})
	}
}

func AddImagesExercise(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	var exercise models.ImagesExercise

	if err := c.BindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	validationErr := validate.Struct(&exercise)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	objectId := primitive.NewObjectID()
	exercise.ID = objectId

	switch language {
	case "de":
		defer cancel()
		result, err := deExercisesCollection.InsertOne(ctx, &exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exercise was not added"})
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	case "kr":
		defer cancel()
		result, err := krExercisesCollection.InsertOne(ctx, &exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exercise was not added"})
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	default:
		fmt.Println(exercise)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect exercise"})
	}
}

func AddAudioExercise(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	var exercise models.AudioExercise

	if err := c.BindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	validationErr := validate.Struct(&exercise)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	objectId := primitive.NewObjectID()
	exercise.ID = objectId

	switch language {
	case "de": // For adding German audio exercise
		defer cancel()
		result, err := deExercisesCollection.InsertOne(ctx, &exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exercise was not added"})
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	case "kr": // For adding Korean audio exercise
		defer cancel()
		result, err := krExercisesCollection.InsertOne(ctx, &exercise)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exercise was not added"})
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	default: // If wrong exercise was provided
		fmt.Println(exercise)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect exercise"})
	}
}

func sendEnDeExercise(ctx context.Context, level int, sendExercise *models.SendTextExercise) error {
	var exercises []models.TextExercise
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 0}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", level}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	cursor, err := tempExercisesCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err = cursor.All(ctx, &exercises); err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, exercise := range exercises {
		sendExercise.ID = exercise.ID
		sendExercise.Type = exercise.Type
		sendExercise.Question = exercise.Question
	}

	return nil
}

func sendDeEnExercise(ctx context.Context, level int, sendExercise *models.SendTextExercise) error {
	var exercises []models.TextExercise
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 1}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", level}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	cursor, err := tempExercisesCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err = cursor.All(ctx, &exercises); err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, exercise := range exercises {
		sendExercise.ID = exercise.ID
		sendExercise.Type = exercise.Type
		sendExercise.Question = exercise.Question
	}

	return nil
}

func sendImageExercise(ctx context.Context, level int, sendExercise *models.SendImageExercise) error {
	var exercises []models.ImageExercise
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 2}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", level}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	cursor, err := tempExercisesCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err = cursor.All(ctx, &exercises); err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, exercise := range exercises {
		sendExercise.ID = exercise.ID
		sendExercise.Type = exercise.Type
		sendExercise.Question = exercise.Question
	}

	return nil
}

func sendImagesExercise(ctx context.Context, level int, sendExercise *models.SendImagesExercise) error {
	var exercises []models.ImagesExercise
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 3}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", level}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	cursor, err := tempExercisesCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err = cursor.All(ctx, &exercises); err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, exercise := range exercises {
		sendExercise.ID = exercise.ID
		sendExercise.Type = exercise.Type
		sendExercise.Word = exercise.Word
		sendExercise.Cards = exercise.Cards
	}

	return nil
}

func sendAudioExercise(ctx context.Context, level int, sendExercise *models.SendAudioExercise) error {
	var exercises []models.AudioExercise
	matchTypeStage := bson.D{{"$match", bson.D{{"type", 4}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", level}}}}
	randomStage := bson.D{{"$sample", bson.D{{"size", 1}}}}

	cursor, err := tempExercisesCollection.Aggregate(ctx, mongo.Pipeline{matchTypeStage, matchLevelStage, randomStage})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err = cursor.All(ctx, &exercises); err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, exercise := range exercises {
		sendExercise.ID = exercise.ID
		sendExercise.Type = exercise.Type
		sendExercise.Question = exercise.Answer
	}

	return nil
}

// generateRandomType is a function that generates number between 0 and 4
// which is a type of exercise
func generateRandomType() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 4

	return rand.Intn(max-min+1) + min
}

// calculateGainExp returns number of experience
// gained by user after solving question
func calculateGainExp(level int) int {
	if level/100 == 0 {
		return 5
	}
	return 1 / (level / (100 - (level / 100))) * ExpMultiplier
}

func checkASRConnection() error {
	req, err := http.NewRequest("GET", "http://localhost:4040/predict", nil)
	if err != nil {
		return err
	}

	res, _ := client.Do(req)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}
	return nil
}
