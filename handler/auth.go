package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) AuthenticateUser(c *gin.Context) {
	var token models.AuthToken

	// Bind JSON to token structure
	if err := c.BindJSON(&token); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&token)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	user, err := h.services.Authorization.AuthenticateUser(token)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
		return
	}

	if user.LastLanguage == "" {
		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}

	lastLanguageProgress, err := h.services.GetLastLanguage(user)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get last language"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "last_language": lastLanguageProgress})
}
