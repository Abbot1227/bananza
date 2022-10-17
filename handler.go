package main

import (
	"Bananza/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	// user endpoint
	user := router.Group("/user")
	{
		user.POST("/login", routes.AuthenticateUser)
	}

	return router
}
