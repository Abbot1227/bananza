package routes

import (
	"Bananza/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

var validate = validator.New()

// AuthenticateUser is a function TODO add description
func AuthenticateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var token models.AuthToken // TODO change maybe

	if err := c.BindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(token)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
}

// ValidateToken is a function TODO add description
func ValidateToken(token string) error {
	oauth2Service, err := oauth2.NewService(context.Background(), option.WithoutAuthentication())
	if err != nil {
		return err
	}

	userInfoService := oauth2.NewUserinfoV2MeService(oauth2Service)

	userInfo, err := userInfoService.Get().Do(googleapi.QueryParameter("access_token", token))
	if err != nil {
		e, _ := err.(*googleapi.Error)
		fmt.Println(e.Message)
		return e
	}
	userInfo.Email
}
