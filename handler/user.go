package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (h *Handler) AddLanguage(c *gin.Context) {
	var inputLanguage models.InputLanguage

	if err := c.BindJSON(&inputLanguage); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&inputLanguage)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	result, err := h.services.User.AddLanguage(inputLanguage)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add language to user"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) UserProfiles(c *gin.Context) {
	profiles, err := h.services.User.FindProfiles()
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get users profiles"})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *Handler) UserProfile(c *gin.Context) {
	user := c.Params.ByName("id")
	userId, _ := primitive.ObjectIDFromHex(user[3:])

	profile, err := h.services.User.FindProfile(userId)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user profile"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UserProgress(c *gin.Context) {
	user := c.Query("id")
	language := c.Query("language")
	userId, _ := primitive.ObjectIDFromHex(user)

	logrus.Println(user + " " + language)

	progress, err := h.services.User.FindProgress(userId, language)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user progress"})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func (h *Handler) UserProgresses(c *gin.Context) {
	user := c.Params.ByName("id")
	userId, _ := primitive.ObjectIDFromHex(user[3:])

	progresses, err := h.services.User.FindProgresses(userId)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user progress"})
		return
	}
	c.JSON(http.StatusOK, progresses)
}

func (h *Handler) UpdateProgress(c *gin.Context) {
	var updateProgress models.UserProgressUpdate

	logrus.Println(updateProgress)

	if err := c.BindJSON(&updateProgress); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(&updateProgress)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
}

func (h *Handler) SetLastLanguage(c *gin.Context) {
	var lastLanguage models.LastLanguageUpdate

	logrus.Println(lastLanguage)

	if err := c.BindJSON(&lastLanguage); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if lastLanguage.ID == "0" && lastLanguage.LastLanguage == "0" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "try again"})
		return
	}

	validationErr := validate.Struct(&lastLanguage)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	userId, _ := primitive.ObjectIDFromHex(lastLanguage.ID)

	if err := h.services.User.SetLastLanguage(userId, lastLanguage.LastLanguage); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update language"})
		return
	}

	languageProgress, err := h.services.User.FindProgress(userId, lastLanguage.LastLanguage)
	if err == mongo.ErrNoDocuments {
		var inputLanguage = models.InputLanguage{Language: lastLanguage.LastLanguage,
			User: lastLanguage.ID}

		_, err = h.services.User.AddLanguage(inputLanguage)
		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		languageProgress, err = h.services.User.FindProgress(userId, lastLanguage.LastLanguage)
		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, languageProgress)
}
