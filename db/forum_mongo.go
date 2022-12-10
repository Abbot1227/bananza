package db

import (
	"Bananza/models"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForumMongo struct {
}

func NewForumMongo() *ForumMongo {
	return &ForumMongo{}
}

func (r *ForumMongo) FindUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error) {
	var user models.User
	filter := bson.D{{"_id", userId}}

	if err := usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	logrus.Println(user)
	return &user, nil
}

func (r *ForumMongo) CreatePost(ctx context.Context, forumPost *models.ForumPost) error {
	result, err := postsCollection.InsertOne(ctx, forumPost)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}

func (r *ForumMongo) GetForumPosts(ctx context.Context) ([]models.ForumPost, error) {
	var forumPosts []models.ForumPost

	cursor, err := postsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &forumPosts); err != nil {
		return nil, err
	}
	logrus.Println(forumPosts)
	return forumPosts, nil
}

func (r *ForumMongo) GetForumPost(ctx context.Context, postId primitive.ObjectID) (*models.ForumPost, error) {
	var forumPost models.ForumPost
	filter := bson.D{{"_id", postId}}

	if err := postsCollection.FindOne(ctx, filter).Decode(&forumPost); err != nil {
		return nil, err
	}
	logrus.Println(forumPost)
	return &forumPost, nil
}

func (r *ForumMongo) CreateComment(ctx context.Context, forumComment *models.ForumComment, postId primitive.ObjectID) error {
	update := bson.D{{"$push", bson.D{{"replies", forumComment}}}}

	result, err := postsCollection.UpdateByID(ctx, postId, update)
	if err != nil {
		return err
	}
	logrus.Println(result)
	return nil
}
