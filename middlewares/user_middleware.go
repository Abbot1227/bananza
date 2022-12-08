package middlewares

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func UserMongoMiddleware() gin.HandlerFunc {
	var usersCollection = db.OpenCollection(db.Client, "users")

	return func(c *gin.Context) {
		var inputLanguage models.InputLanguage

		if err := c.BindJSON(&inputLanguage); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		user, _ := primitive.ObjectIDFromHex(inputLanguage.User)
		filter := bson.D{{"_id", user}}

		if err := usersCollection.FindOne(context.Background(), filter).Decode(nil); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		}

		c.Next()
	}
}
