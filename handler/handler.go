package handler

import (
	"Bananza/routes"
	"Bananza/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

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

	// Endpoints for user authentication and progress
	user := router.Group("/users")
	{
		user.POST("/login", routes.AuthenticateUser)
	}

	return router
}
