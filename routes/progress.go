package routes

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

var userProgressCollection = db.OpenCollection(db.Client, "userProgress")

// AddLanguage is a function TODO add description
func AddLanguage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var inputLanguage models.InputLanguage // TODO change maybe

	if err := c.BindJSON(&inputLanguage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&inputLanguage)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	user, _ := primitive.ObjectIDFromHex(inputLanguage.User)
	filter := bson.D{{"_id", user}}

	// TODO move to middleware
	if err := usersCollection.FindOne(ctx, filter).Decode(nil); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		fmt.Println(err)
		return
	}

	userProgress := models.UserProgress{ID: primitive.NewObjectID(),
		Language: inputLanguage.Language,
		Progress: 0,
		User:     user}

	// Inserting new language into database
	result, insertErr := userProgressCollection.InsertOne(ctx, &userProgress)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "language was not added"})
		fmt.Println(insertErr)
		return
	}

	defer cancel()

	// Return result of existing user
	c.JSON(http.StatusOK, result)
	return

}

// UserProgress is a function
func UserProgress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	userID := c.Query("id")
	language := c.Query("language")
	user, _ := primitive.ObjectIDFromHex(userID)

	fmt.Println(userID + " " + language)

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user", user}},
				bson.D{{"language", language}},
			},
		},
	}
	var languageProgress models.UserProgress

	if err := userProgressCollection.FindOne(ctx, filter).Decode(&languageProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()

	fmt.Println(languageProgress)

	c.JSON(http.StatusOK, languageProgress)
}

// UserProgresses is a function TODO add description
func UserProgresses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID[3:])
	filter := bson.D{{"user", docID}}

	var userProgress []models.UserProgress

	cursor, err := userProgressCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	if err = cursor.All(ctx, &userProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}
	defer cancel()

	fmt.Println(userProgress)

	c.JSON(http.StatusOK, userProgress)
}

// UpdateProgress is a function
func UpdateProgress(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
}

// Update progress
// Get audio file from user endpoint
// Do refactoring
// Set exercises model
