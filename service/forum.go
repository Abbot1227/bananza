package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ForumService struct {
	repo db.Forum
}

func NewForumService(repo db.Forum) *ForumService {
	return &ForumService{repo: repo}
}

func (s *ForumService) AddPost(inputForumPost *models.InputForumPost, forumPost *models.ForumPost) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Fetching user document from db to assign Name to Author of forum post
	userId, _ := primitive.ObjectIDFromHex(inputForumPost.UserId)
	user, err := s.repo.FindUser(ctx, userId)
	if err != nil {
		return err
	}
	defer cancel()

	// Creating values for forum post
	postId := primitive.NewObjectID()
	author := user.Name
	createdAt := primitive.NewDateTimeFromTime(time.Now())
	replies := []models.ForumComment{} // Initially there are no replies to post

	// Assigning values to forum post
	forumPost.ID = postId
	forumPost.Title = inputForumPost.Title
	forumPost.Text = inputForumPost.Text
	forumPost.Author = author
	forumPost.CreatedAt = createdAt
	forumPost.Replies = replies

	if err := s.repo.CreatePost(ctx, forumPost); err != nil {
		return nil
	}

	return nil
}

func (s *ForumService) GetForumTitles(skip int) (*[]models.SendForumTitles, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var forumTitles []models.SendForumTitles

	forumPosts, err := s.repo.GetForumPosts(ctx, skip)
	if err != nil {
		return nil, err
	}
	defer cancel()

	for _, post := range forumPosts {
		forumTitles = append(forumTitles, models.SendForumTitles{
			ID:        post.ID,
			Title:     post.Title,
			CreatedAt: post.CreatedAt,
		})
	}

	return &forumTitles, nil
}

func (s *ForumService) GetForumPost(postId primitive.ObjectID) (*models.ForumPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	forumPost, err := s.repo.GetForumPost(ctx, postId)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return forumPost, nil
}

func (s *ForumService) AddComment(inputComment *models.InputForumComment, postComment *models.ForumComment, postId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Fetching user document from db to assign Name to Author of comment post
	userId, _ := primitive.ObjectIDFromHex(inputComment.UserId)
	user, err := s.repo.FindUser(ctx, userId)
	if err != nil {
		return err
	}
	defer cancel()

	// Creating values for post comment
	commentId := primitive.NewObjectID()
	author := user.Name
	createdAt := primitive.NewDateTimeFromTime(time.Now())

	// Assigning values to post comment
	postComment.ID = commentId
	postComment.Text = inputComment.Text
	postComment.Author = author
	postComment.CreatedAt = createdAt

	if err := s.repo.CreateComment(ctx, postComment, postId); err != nil {
		return err
	}

	return nil
}

func (s *ForumService) RemovePost(postId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := s.repo.DeletePost(ctx, postId); err != nil {
		return err
	}
	defer cancel()

	return nil
}
