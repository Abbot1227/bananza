package handler

import (
	"Bananza/middlewares"
	"Bananza/service"
	"net/http"

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
		user.PUT("/progress", h.UpdateProgress)  // Done Test add middlewares check if exists
		user.PUT("/lastlang", h.SetLastLanguage) // Done
		// D
		user.DELETE("/:id", h.RemoveUser) // Done
	}

	exercise := router.Group("/exercises")
	{
		// C
		exercise.POST("/new", h.SendExercise)                     // Done
		exercise.POST("/answer", h.SendAnswer)                    // Done
		exercise.POST("/audio/:lang", h.LoadAudio)                // Done
		exercise.POST("add/teximg/:lang", h.AddTextImageExercise) // Done
		exercise.POST("add/imgs/:lang", h.AddImagesExercise)      // Done
		exercise.POST("add/audio/:lang", h.AddAudioExercise)      // Done
		// R
		exercise.GET("/mul", h.SetMultiplier) // Done
		exercise.GET("/asr", h.SetASRUrl)     // Done
	}

	forum := router.Group("/forum")
	{
		// C
		forum.POST("/post", h.CreatePost) // Done
		// R
		forum.GET("/", h.ForumTitles)                                           // Add pagination
		forum.GET("/post/:id", h.ForumPost).Use(middlewares.PostIdMiddleware()) // Done
		// U
		forum.PUT("/post/:id", h.AddComment).Use(middlewares.PostIdMiddleware()) // Done
		// D
		forum.DELETE("/post/:id", h.RemovePost).Use(middlewares.PostIdMiddleware()) // Done
	}

	shop := router.Group("/shop")
	{
		// C
		shop.POST("/buy", h.BuyAvatar) // Done
		// R
		shop.GET("/", h.Avatars) // Done
		// U
		shop.PUT("/set", h.SetAvatar) // Done
	}

	grammar := router.Group("/grammar")
	{
		// R
		grammar.GET("/")
		grammar.GET("/dictionary")
	}

	return router
}

func checkHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "is working"})
}
