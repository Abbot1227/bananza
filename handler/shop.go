package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) BuyAvatar(c *gin.Context) {
	var inputAvatarPurchase models.InputAvatarPurchase

	if err := c.BindJSON(&inputAvatarPurchase); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputAvatarPurchase)

	// Ensure that data we receive is correct
	if validationErr := validate.Struct(&inputAvatarPurchase); validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err := h.services.Shop.BuyAvatar(&inputAvatarPurchase); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar purchased successfully"})
}

func (h *Handler) Avatars(c *gin.Context) {
	avatars, err := h.services.Shop.GetAvatars()
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, avatars)
}

func (h *Handler) SetAvatar(c *gin.Context) {
	var inputAvatarSet models.InputAvatarSet

	if err := c.BindJSON(&inputAvatarSet); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputAvatarSet)

	// Ensure that data we receive is correct
	if validationErr := validate.Struct(&inputAvatarSet); validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err := h.services.Shop.SetAvatar(&inputAvatarSet); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar set successfully"})
}
