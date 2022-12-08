package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

}

func (h *Handler) UserProgresses(c *gin.Context) {

}

func (h *Handler) UpdateProgress(c *gin.Context) {

}

func (h *Handler) SetLastLanguage(c *gin.Context) {

}
