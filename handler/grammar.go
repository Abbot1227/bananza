package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) Grammar(c *gin.Context) {
	var inputGrammar models.InputDictionary

	if err := c.BindJSON(&inputGrammar); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputGrammar)

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&inputGrammar)
	if validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	inputGrammar.Level = inputGrammar.Level / 100

	grammar, err := h.services.Grammar.GetGrammar(inputGrammar)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(grammar)

	c.JSON(http.StatusOK, grammar)
}

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

	inputDictionary.Level = inputDictionary.Level / 100

	dictionary, err := h.services.Grammar.GetDictionary(inputDictionary)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(dictionary)

	c.JSON(http.StatusOK, dictionary)
}
