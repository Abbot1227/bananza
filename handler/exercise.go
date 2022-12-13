package handler

import (
	"Bananza/models"
	"Bananza/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var ExpMultiplier = 50

// SendExercise godoc
// @Summary
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

// SendAnswer godoc
// @Summary
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

	answer, err := h.services.Exercise.GetRightAnswer(inputAnswer.ID)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.Println("Right answer: ", answer)

	expToAdd := calculateGainExp(inputAnswer.Level)

	if inputAnswer.Answer == answer {
		c.JSON(http.StatusOK, gin.H{"correct": "true", "answer": answer, "exp": expToAdd})
	} else {
		c.JSON(http.StatusOK, gin.H{"correct": "false", "answer": answer, "exp": 0})
		return
	}

	if err := h.services.Exercise.UpdateProgress(inputAnswer.LanguageId, expToAdd); err != nil {
		logrus.Error(err.Error())
		logrus.Println("could not update user's progress")
	}
}

// LoadAudio godoc
// @Summary
func (h *Handler) LoadAudio(c *gin.Context) {
	// The FormFile function takes in the POST input id file
	c.Request.ParseMultipartForm(64 << 20)

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
	logrus.Println(questionId, languageId, level)

	// Get user answer in text format
	userAnswer, err := h.services.Exercise.GetAudioAnswer(file, language)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get right answer for user question
	rightAnswer, err := h.services.Exercise.GetRightAnswer(questionId[0])
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Println(userAnswer, " ", rightAnswer)

	exp, _ := strconv.Atoi(level[0])
	expToAdd := calculateGainExp(exp)

	if userAnswer == rightAnswer {
		c.JSON(http.StatusOK, gin.H{"correct": "true", "answer": rightAnswer, "exp": expToAdd, "userAnswer": userAnswer})
	} else {
		c.JSON(http.StatusOK, gin.H{"correct": "false", "answer": rightAnswer, "exp": 0, "userAnswer": userAnswer})
		return
	}

	if err := h.services.Exercise.UpdateProgress(languageId[0], expToAdd); err != nil {
		logrus.Error(err.Error())
		logrus.Println("could not update user's progress")
	}
}

// SetMultiplier godoc
// @Summary
func (h *Handler) SetMultiplier(c *gin.Context) {
	multiplier := c.Query("set")
	newMultiplier, _ := strconv.Atoi(multiplier)

	ExpMultiplier = newMultiplier
	c.JSON(http.StatusOK, gin.H{})
}

// SetASRUrl godoc
// @Summary
func (h *Handler) SetASRUrl(c *gin.Context) {
	url := c.Query("set")

	service.ASRUrl = url
	c.JSON(http.StatusOK, gin.H{})
}

// calculateGainExp returns number of experience
// gained by user after solving question
func calculateGainExp(level int) float64 {
	if level/100 == 0 || level/100 == 1 {
		return 5
	}
	return (100 / float64(level)) * float64(ExpMultiplier)
}
