package main

import (
	"Bananza/db"
	"Bananza/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
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
	db.Client = db.ConnectDB()
	srv := new(Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatal(err)
	}

	router = gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	// user endpoint
	user := router.Group("/user")
	{
		user.POST("/login", routes.AuthenticateUser)
	}

	router.Run(":8080")
}
