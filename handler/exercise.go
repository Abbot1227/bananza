package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var ExpMultiplier = 15

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
	var inputAnswer models.InputAnswer

	if err := c.BindJSON(&inputAnswer); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputAnswer)

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&inputAnswer)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	answer, err := h.services.Exercise.GetRightAnswer(inputAnswer.Answer)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expToAdd := calculateGainExp(inputAnswer.Level)

	if inputAnswer.Answer == answer {
		c.JSON(http.StatusOK, gin.H{"correct": "true", "answer": answer, "exp": expToAdd})
	} else {
		c.JSON(http.StatusOK, gin.H{"correct": "false", "answer": answer, "exp": 0})
	}

	if err := h.services.Exercise.UpdateProgress(inputAnswer.LanguageId, expToAdd); err != nil {
		logrus.Error(err.Error())
		logrus.Println("could not update user's progress")
	}
}

func (h *Handler) LoadAudio(c *gin.Context) {
	// The FormFile function takes in the POST input id file
	c.Request.ParseMultipartForm(32 << 20)

	languageParam := c.Params.ByName("lang")
	language := languageParam[5:]

	questionId := c.Request.MultipartForm.Value["id"]
	languageId := c.Request.MultipartForm.Value["languageId"]
	level := c.Request.MultipartForm.Value["level"]
	file, _, err := c.Request.FormFile("mp3")
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// TODO Остановился здесь
	answer, err := h.services.Exercise.GetRightAnswer(inputAnswer.Answer)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) SetMultiplier(c *gin.Context) {
	multiplier := c.Query("set")
	newMultiplier, _ := strconv.Atoi(multiplier)

	ExpMultiplier = newMultiplier
	c.JSON(http.StatusOK, gin.H{})
}

// calculateGainExp returns number of experience
// gained by user after solving question
func calculateGainExp(level int) int {
	if level/100 == 0 {
		return 5
	}
	return 1 / (level / (100 - (level / 100))) * ExpMultiplier
}
