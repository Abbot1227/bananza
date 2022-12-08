package handler

import (
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

	router.GET("/health", checkHealth)

	user := router.Group("/users")
	{
		// C
		user.POST("/login", h.AuthenticateUser) // Done
		user.POST("/progress", h.AddLanguage)   // Done add middlewares check if language exists and user exists or not
		// R
		user.GET("/", h.UserProfiles)               // Done add middlewares authorization only admin
		user.GET("/:id", h.UserProfile)             // Done authorization only admin
		user.GET("/progresslang", h.UserProgress)   // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		user.GET("/progress/:id", h.UserProgresses) // Done add middlewares authorization only admin проверить что только юзер с тем id может запрашивать свои
		// U
		user.PUT("/progress", h.UpdateProgress)  // Test add middlewares check if exists
		user.PUT("/lastlang", h.SetLastLanguage) // Done
		// D
	}

	exercise := router.Group("/exercises")
	{
		// C
		exercise.POST("/new", h.SendExercise)
		exercise.POST("/answer", h.SendAnswer)
		exercise.POST("/audio/:lang", h.LoadAudio)
		exercise.POST("add/teximg/:lang", h.AddTextImageExercise)
		exercise.POST("add/imgs/:lang", h.AddImagesExercise)
		exercise.POST("add/audio/:lang", h.AddAudioExercise)
		// R
		exercise.GET("/mul", h.SetMultiplier)
	}

	return router
}

func checkHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "is working"})
}
