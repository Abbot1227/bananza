package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) Dictionary(c *gin.Context) {
	var inputDictionary models.InputDictionary

	if err := c.BindJSON(&inputDictionary); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputDictionary)

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&inputDictionary)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	dictionary, err := h.services.Grammar.GetDictionary(inputDictionary)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(dictionary)

	c.JSON(http.StatusOK, dictionary)
}
