package main

import (
	"Bananza/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// Testing docker purpose only
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "docker test"})
	})

	// user endpoint
	user := router.Group("/users")
	{
		// C
		user.POST("/login", routes.AuthenticateUser) // Done
		user.POST("/progress", routes.AddLanguage)   // Done add middlewares check if language exists and user exists or not
		// R
		user.GET("/", routes.UserProfiles)             // Done add middlewares authorization only admin
		user.GET("/:id", routes.UserProfile)           // Done authorization only admin
		user.GET("/progresslang", routes.UserProgress) // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		user.GET("/progress", routes.UserProgresses)
		// U
		user.PUT("/progress", routes.UpdateProgress) // Test add middlewares check if exists
		user.PUT("/lastlang", routes.SetLastLanguage)
		// TODO delete progress and-or user
	}

	exercise := router.Group("/exercises")
	{
		exercise.POST("/mic", routes.LoadAudio) // Test
		exercise.GET("/:num", routes.NextExercise)
	}

	router.Run(":8080")
	//router.RunTLS(":8080", "cert.pem", "key.pem")
}

// Проверять соединение с ML сервером, убирать задания на speech если недоступно

// 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
