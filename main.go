package main

import (
	"Bananza/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		user.POST("/progress", routes.AddLanguage)   // Test add middlewares check if language exists and user exists or not
		// R
		user.GET("/", routes.UserProfiles)         // Done add middlewares authorization only admin
		user.GET("/:id", routes.UserProfile)       // Done authorization only admin
		user.GET("/progress", routes.UserProgress) // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		// U
		user.PUT("/progress", routes.UpdateProgress) // Test add middlewares check if exists
	}

	router.Run(":8080")
}

// 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
