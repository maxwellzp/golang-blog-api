package service

import (
	"context"
	"database/sql"
	"errors"
	"maxwellzp/blog-api/internal/model"
	"maxwellzp/blog-api/internal/repository"
	"strings"
)

type CommentService interface {
	Create(ctx context.Context, userID, blogID int64, content string) (*model.Comment, error)
	GetByID(ctx context.Context, id int64) (*model.Comment, error)
	Update(ctx context.Context, id int64, content string) error
	Delete(ctx context.Context, id int64) error
	ListByBlogID(ctx context.Context, blogID int64, limit, offset int) ([]*model.Comment, error)
	IsOwner(ctx context.Context, commentID, userID int64) (bool, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) Create(ctx context.Context, userID, blogID int64, content string) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content is empty")
	}

	comment := &model.Comment{
		UserID:  userID,
		BlogID:  blogID,
		Content: content,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentService) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *commentService) Update(ctx context.Context, id int64, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("content cannot be empty")
	}

	comment := &model.Comment{
		ID:      id,
		Content: content,
	}
	return s.repo.Update(ctx, comment)
}

func (s *commentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *commentService) ListByBlogID(ctx context.Context, blogID int64, limit, offset int) ([]*model.Comment, error) {
	return s.repo.ListByBlogID(ctx, blogID, limit, offset)
}

var ErrCommentNotFound = errors.New("comment not found")

func (s *commentService) IsOwner(ctx context.Context, commentID, userID int64) (bool, error) {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrCommentNotFound
		}
		return false, err
	}
	return comment.UserID == userID, nil
}
