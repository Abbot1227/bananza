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
}
