package service

import (
	"context"
	"database/sql"
	"errors"
	"maxwellzp/blog-api/internal/model"
	"maxwellzp/blog-api/internal/repository"
	"strings"
)

type BlogService interface {
	Create(ctx context.Context, userId int64, title, content string) (*model.Blog, error)
	GetByID(ctx context.Context, id int64) (*model.Blog, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, title, content string) error
	List(ctx context.Context) ([]*model.Blog, error)
	IsOwner(ctx context.Context, blogID, userID int64) (bool, error)
}

type blogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) BlogService {
	return &blogService{repo: repo}
}

func (s *blogService) Create(ctx context.Context, userId int64, title, content string) (*model.Blog, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return nil, errors.New("title or content is empty")
	}

	blog := &model.Blog{
		UserID:  userId,
		Title:   title,
		Content: content,
	}

	if err := s.repo.Create(ctx, blog); err != nil {
		return nil, err
	}
	return blog, nil
}

func (s *blogService) GetByID(ctx context.Context, id int64) (*model.Blog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *blogService) Update(ctx context.Context, id int64, title, content string) error {
	if title == "" || content == "" {
		return errors.New("title and content cannot be empty")
	}

	blog := &model.Blog{
		ID:      id,
		Title:   title,
		Content: content,
	}
	return s.repo.Update(ctx, blog)
}

func (s *blogService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *blogService) List(ctx context.Context) ([]*model.Blog, error) {
	return s.repo.List(ctx)
}

var ErrBlogNotFound = errors.New("blog not found")

func (s *blogService) IsOwner(ctx context.Context, blogID, userID int64) (bool, error) {
	blog, err := s.repo.GetByID(ctx, blogID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrBlogNotFound
		}
		return false, err
	}
	return blog.UserID == userID, nil
}
