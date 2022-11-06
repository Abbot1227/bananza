package main

import (
	"Bananza/routes"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// @title TODO App API
// @version 1.0
// description API server for todo application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var router *gin.Engine

func main() {
	/*srv := new(Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatal(err)
	}*/

	router = gin.New()
	router.Use(gin.Logger())
	// CORS configuration
	router.Use(cors.Default())

	// user endpoint
	user := router.Group("/users")
	{
		// C
		user.POST("/login", routes.AuthenticateUser) // Done
		user.POST("/progress", routes.AddLanguage)   // Done add middlewares check if language exists and user exists or not
		// R
		user.GET("/", routes.UserProfiles)         // Done add middlewares authorization only admin
		user.GET("/:id", routes.UserProfile)       // Done authorization only admin
		user.GET("/progress", routes.UserProgress) // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		// U
		user.PUT("/progress", routes.UpdateProgress) // Test add middlewares check if exists
		// TODO delete progress and-or user
	}

	exercise := router.Group("/exercises")
	{
		exercise.POST("/mic", routes.MicroShit) // Test
	}

	r2 := gin.Default()
	r2.POST("/predict", func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		// the FormFile function takes in the POST input id file
		file, header, err := c.Request.FormFile("uploaded_file")
		if err != nil {
			fmt.Println("Error when requesting file nigger: " + err.Error())
			return
		}
		defer file.Close()
		err = c.SaveUploadedFile(header, "D:\\Programming\\Golang\\github.com\\abbot1227\\bananza\\"+header.Filename)
		if err != nil {
			fmt.Println(err)
		}

		defer cancel()

		c.JSON(http.StatusOK, gin.H{"text": "sample"})
	})

	//go r2.Run(":4040")

	router.Run(":8080")
}

// Проверять соединение с ML сервером, убирать задания на speech если недоступно

// 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
