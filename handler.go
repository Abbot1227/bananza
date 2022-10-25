package main

import (
	"Bananza/routes"
	"Bananza/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

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
