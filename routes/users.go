package routes

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
	"io/ioutil"
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

	var user models.User

	if err := validateToken(ctx, token.Token, &user); err != nil {
		// If there is no user create new one
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		}

		// Inserting new user into database
		result, insertErr := usersCollection.InsertOne(ctx, &user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			fmt.Println(insertErr)
			return
		}

		// Return result of existing user
		c.JSON(http.StatusOK, result)
		return
	}
	defer cancel()

	// TODO вынести в отдельную функцию
	filter := bson.D{{"user", user.ID}}
	var userProgress []models.UserProgress

	cursor, err := userProgressCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &userProgress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	// Return user info
	c.JSON(http.StatusOK, gin.H{"user": user, "language": userProgress})
}

// UserProfiles is a function TODO add description
func UserProfiles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var users []models.User

	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
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
		return
	}
	defer cancel()

	fmt.Println(user)

	c.JSON(http.StatusOK, user)
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

// Подумать над интеграцией с Мирасом
// Сделать рефакторинг

// Запасной вариант

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func getGoogleUserInfo(token string) (*UserInfo, error) {
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
