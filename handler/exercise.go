package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) SendExercise(c *gin.Context) {
	var acquireExercise models.AcquireExercise

	if err := c.BindJSON(&acquireExercise); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(acquireExercise)

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&acquireExercise)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	exerciseType, err := h.services.Exercise.GetExerciseType()
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println("Exercise type: ", exerciseType)

	// Reducing experience amount to single value
	acquireExercise.Exp = acquireExercise.Exp / 100

	switch exerciseType {
	case 0: // English to Chosen language exercise
		var exercise models.SendTextExercise

		if err := h.services.Exercise.GetEnLnExercise(acquireExercise, &exercise); err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		logrus.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	case 1: // Chosen language to English
		var exercise models.SendTextExercise

		if err := h.services.Exercise.GetLnEnExercise(acquireExercise, &exercise); err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		logrus.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	case 2: // Translate one image
		var exercise models.SendImageExercise

		if err := h.services.Exercise.GetImageExercise(acquireExercise, &exercise); err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		logrus.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	case 3: // Choose correctly translated image
		var exercise models.SendImagesExercise

		if err := h.services.Exercise.GetImagesExercise(acquireExercise, &exercise); err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		logrus.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	case 4: // Correctly spell word
		var exercise models.SendAudioExercise

		if err := h.services.Exercise.GetAudioExercise(acquireExercise, &exercise); err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		logrus.Println(exercise)

		c.JSON(http.StatusOK, exercise)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}

func (h *Handler) SendAnswer(c *gin.Context) {

}

func (h *Handler) LoadAudio(c *gin.Context) {

}

func (h *Handler) AddTextImageExercise(c *gin.Context) {

}

func (h *Handler) AddImagesExercise(c *gin.Context) {

}

func (h *Handler) AddAudioExercise(c *gin.Context) {

}

func (h *Handler) SetMultiplier(c *gin.Context) {

}
