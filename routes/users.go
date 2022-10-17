package routes

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"
)

var validate = validator.New()
var usersCollection = db.OpenCollection(db.Client, "users")

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

	if err := validateToken(ctx, token.Token); err != nil {
		// If there is no user create new one
		if err == mongo.ErrNoDocuments {
			// Getting user google account info from Google api
			userInfo, err := getUserInfo(token.Token)
			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			// Вынести модель юзера в общую зону и внутри валидации токена добавлять данные в неё
			user := models.User{ID: primitive.NewObjectID(),
				Email:     userInfo.Email,
				Name:      userInfo.Name,
				FirstName: userInfo.GivenName,
				LastName:  userInfo.FamilyName,
				UserId:    userInfo.Id,
				AvatarURL: userInfo.Picture,
			}

			// Inserting new user into database
			result, insertErr := usersCollection.InsertOne(ctx, &user)
			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
				fmt.Println(insertErr)
				return
			}

			c.JSON(http.StatusOK, result)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cancel()

	// Return message "user exists" if user exists
	c.JSON(http.StatusOK, gin.H{"user-status": "already exists"})
}

// Просто проверка работы базы данных
// UserProfile is a function TODO returns users information
func UserProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	filter := bson.D{{"userId", "1234567810"}}
	var user models.User

	if err := usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "the requested user was not found"})
		return
	}
	defer cancel()

	// Return user with given user id
	c.JSON(http.StatusOK, user)
}

// validateToken is a function TODO add description
func validateToken(ctx context.Context, token string) error {
	userInfo, err := getUserInfo(token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	userId := userInfo.Id

	// Filter to find user with specified Google ID in users collection
	filter := bson.D{{"userId", userId}}
	var user bson.M // To store user's info in mongodb object

	// Obtain user info from users collection and store it in user object
	// if not found return error otherwise return nil
	if err = usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		log.Fatal(err)
	}

	return nil
}

// getUserInfo is a function TODO add description
func getUserInfo(token string) (*oauth2.Userinfoplus, error) {
	oauth2Service, err := oauth2.NewService(context.Background(), option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}
	userInfoService := oauth2.NewUserinfoV2MeService(oauth2Service)

	// Getting user google account info from Google api
	userInfo, err := userInfoService.Get().Do(googleapi.QueryParameter("access_token", token))
	if err != nil {
		e, _ := err.(*googleapi.Error)
		return nil, e
	}
	return userInfo, nil
}
