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

	var inputLanguage models.InputLanguage

	if err := c.BindJSON(&inputLanguage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(inputLanguage)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	userId, _ := primitive.ObjectIDFromHex(inputLanguage.User)

	userProgress := &models.UserProgress{ID: primitive.NewObjectID(),
		Language: inputLanguage.Language,
		Level:    0,
		User:     userId,
	}

	result, insertErr := userProgressCollection.InsertOne(ctx, &userProgress)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "language was not added"})
		fmt.Println(insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

// UserProgress is a function
func UserProgress(c *gin.Context) {
	userID := c.Params.ByName("id")
	user, _ := primitive.ObjectIDFromHex(userID[3:])

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	filter := bson.D{{"user", user}}
	var userProgress []models.UserProgress

	cursor, err := userProgressCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &userProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()

	fmt.Println(user)

	c.JSON(http.StatusOK, userProgress)
}

// UpdateProgress is a function
func UpdateProgress(c *gin.Context) {

}
