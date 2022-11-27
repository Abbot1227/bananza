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

	var token models.AuthToken

	if err := c.BindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&token)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	var user models.User // TODO change maybe

	if err := validateToken(ctx, token.Token, &user); err != nil {
		// If there is no user create new one
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		// Getting user google account info from Google api
		userInfo, err := getUserInfo(token.Token)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// Assigning userInfo fields to user fields
		{
			user.ID = primitive.NewObjectID()
			user.Email = userInfo.Email
			user.Name = userInfo.Name
			user.FirstName = userInfo.GivenName
			user.LastName = userInfo.FamilyName
			user.UserId = userInfo.Id
			user.AvatarURL = userInfo.Picture
			user.LastLanguage = "" // TODO change
		}

		// Inserting new user into database
		result, insertErr := usersCollection.InsertOne(ctx, &user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			fmt.Println(insertErr)
			defer cancel()
			return
		}
		defer cancel()

		// Return result of existing user
		c.JSON(http.StatusOK, result)
		return
	}
	defer cancel()

	if user.LastLanguage == "" {
		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}

	// TODO вынести в отдельную функцию
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user", user.ID}},
				bson.D{{"language", user.LastLanguage}},
			},
		},
	}
	var lastLanguageProgress models.UserProgress

	if err := userProgressCollection.FindOne(ctx, filter).Decode(&lastLanguageProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Return user info
	c.JSON(http.StatusOK, gin.H{"user": user, "last_language": lastLanguageProgress})
}

// UserProfiles is a function TODO add description
func UserProfiles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var users []models.User

	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}

	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}
	defer cancel()

	fmt.Println(users)

	c.JSON(http.StatusOK, users)
}

// UserProfile is a function TODO returns users information
func UserProfile(c *gin.Context) {
	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID[3:])

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	filter := bson.D{{"_id", docID}}
	var user models.User

	if err := usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}
	defer cancel()

	fmt.Println(user)

	c.JSON(http.StatusOK, user)
}

func SetLastLanguage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var lastLanguage models.LastLanguageUpdate

	if err := c.BindJSON(&lastLanguage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		defer cancel()
		return
	}
	fmt.Println(lastLanguage)

	// Ensure that data we receive is correct
	validationErr := validate.Struct(&lastLanguage)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	update := bson.D{
		{"$set",
			bson.D{
				{"lastlanguage", lastLanguage.LastLanguage},
			},
		},
	}
	userID, _ := primitive.ObjectIDFromHex(lastLanguage.ID)

	_, err := usersCollection.UpdateByID(ctx, userID, update)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		defer cancel()
	}

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"user", userID}},
				bson.D{{"language", lastLanguage.LastLanguage}},
			},
		},
	}
	var languageProgress models.UserProgress

	if err := userProgressCollection.FindOne(ctx, filter).Decode(&languageProgress); err != nil {
		if err == mongo.ErrNoDocuments {
			newUserProgress := models.UserProgress{ID: primitive.NewObjectID(),
				Language: lastLanguage.LastLanguage,
				Level:    0,
				User:     userID}

			insertResult, insertErr := userProgressCollection.InsertOne(ctx, &newUserProgress)
			if insertErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				defer cancel()
			}
			defer cancel()
			fmt.Println(insertResult)

			c.JSON(http.StatusOK, newUserProgress)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
	}
	defer cancel()

	c.JSON(http.StatusOK, languageProgress)
}

// validateToken is a function TODO add description
func validateToken(ctx context.Context, token string, user *models.User) error {
	userInfo, err := getUserInfo(token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	userId := userInfo.Id

	// Filter to find user with specified Google ID in users collection
	filter := bson.D{{"userid", userId}}

	// Obtain user info from users collection and store it in user object
	// if not found return error otherwise return nil
	if err = usersCollection.FindOne(ctx, filter).Decode(user); err != nil {
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
