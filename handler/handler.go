package handler

import (
	"Bananza/routes"
	"Bananza/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "is working"})
	})

	// Endpoints for user authentication and progress
	user := router.Group("/users")
	{
		// C
		user.POST("/login", routes.AuthenticateUser) // Done
		user.POST("/progress", routes.AddLanguage)   // Done add middlewares check if language exists and user exists or not
		// R
		user.GET("/", routes.UserProfiles)               // Done add middlewares authorization only admin
		user.GET("/:id", routes.UserProfile)             // Done authorization only admin
		user.GET("/progresslang", routes.UserProgress)   // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		user.GET("/progress/:id", routes.UserProgresses) // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		// U
		user.PUT("/progress", routes.UpdateProgress) // Test add middlewares check if exists
		user.PUT("/lastlang", routes.SetLastLanguage)
		// TODO delete progress and-or user
	}

	return router
}
