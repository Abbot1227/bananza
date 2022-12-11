package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// AddTextImageExercise godoc
// @Summary
func (h *Handler) AddTextImageExercise(c *gin.Context) {
	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	if language != "de" && language != "kr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong language"})
	}

	var exercise models.TextExercise
	logrus.Println(exercise)

	if err := c.BindJSON(&exercise); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&exercise)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err := h.services.Exercise.CreateTextImageExercise(exercise, language); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

// AddImagesExercise godoc
// @Summary
func (h *Handler) AddImagesExercise(c *gin.Context) {
	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	if language != "de" && language != "kr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong language"})
	}

	var exercise models.ImagesExercise
	logrus.Println(exercise)

	if err := c.BindJSON(&exercise); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&exercise)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err := h.services.Exercise.CreateImagesExercise(exercise, language); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

// AddAudioExercise godoc
// @Summary
func (h *Handler) AddAudioExercise(c *gin.Context) {
	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	if language != "de" || language != "kr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong language"})
	}

	var exercise models.AudioExercise
	logrus.Println(exercise)

	if err := c.BindJSON(&exercise); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&exercise)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err := h.services.Exercise.CreateAudioExercise(exercise, language); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}
