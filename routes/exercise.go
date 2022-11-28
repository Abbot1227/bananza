package routes

import (
	"Bananza/db"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

var tempExercisesCollection = db.OpenCollection(db.Client, "tempExercises")

func NextExercise(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	switch {

	}
}
