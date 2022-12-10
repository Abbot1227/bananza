package handler

import (
	"Bananza/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var inputForumPost models.InputForumPost

	if err := c.BindJSON(&inputForumPost); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputForumPost)

	// Ensure that data we receive is correct
	if validationErr := validate.Struct(&inputForumPost); validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var forumPost models.ForumPost

	if err := h.services.Forum.AddPost(&inputForumPost, &forumPost); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(forumPost)

	c.JSON(http.StatusOK, forumPost)
}

func (h *Handler) ForumTitles(c *gin.Context) {
	forumTitles, err := h.services.Forum.GetForumTitles()
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(forumTitles)

	c.JSON(http.StatusOK, forumTitles)
}

func (h *Handler) ForumPost(c *gin.Context) {
	post := c.Params.ByName("id")
	postId, _ := primitive.ObjectIDFromHex(post[3:])

	forumPost, err := h.services.Forum.GetForumPost(postId)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, forumPost)
}

func (h *Handler) AddComment(c *gin.Context) {
	post := c.Params.ByName("id")
	postId, _ := primitive.ObjectIDFromHex(post[3:])

	var inputComment models.InputForumComment

	if err := c.BindJSON(&inputComment); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(inputComment)

	// Ensure that data we receive is correct
	if validationErr := validate.Struct(&inputComment); validationErr != nil {
		logrus.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var postComment models.ForumComment

	if err := h.services.Forum.AddComment(&inputComment, &postComment, postId); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.Println(postComment)

	c.JSON(http.StatusOK, postComment)
}

func (h *Handler) RemovePost(c *gin.Context) {
	post := c.Params.ByName("id")
	postId, _ := primitive.ObjectIDFromHex(post[3:])

	if err := h.services.Forum.RemovePost(postId); err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
