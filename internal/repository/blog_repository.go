package repository

import (
	"context"
	"database/sql"
	"maxwellzp/blog-api/internal/model"
	"time"
)

type BlogRepository interface {
	Create(ctx context.Context, blog *model.Blog) error
	GetByID(ctx context.Context, id int64) (*model.Blog, error)
	Update(ctx context.Context, blog *model.Blog) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*model.Blog, error)
}

type blogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Create(ctx context.Context, blog *model.Blog) error {
	query := "INSERT INTO blog (user_id, title, content) VALUES(?, ?, ?)"

	res, err := r.db.ExecContext(ctx, query, blog.UserID, blog.Title, blog.Content)
	if err != nil {
		return err
	}

	blog.ID, err = res.LastInsertId()
	return err
}

func (r *blogRepository) GetByID(ctx context.Context, id int64) (*model.Blog, error) {
	query := "SELECT id, user_id, title, content FROM blog WHERE id = ? AND deleted_at IS NULL"

	row := r.db.QueryRowContext(ctx, query, id)

	blog := &model.Blog{}
	if err := row.Scan(&blog.ID, &blog.UserID, &blog.Title, &blog.Content); err != nil {
		return nil, err
	}
	return blog, nil
}

func (r *blogRepository) Update(ctx context.Context, blog *model.Blog) error {
	query := "UPDATE blog SET title = ?, content = ? WHERE id = ? AND deleted_at IS NULL"

	_, err := r.db.ExecContext(ctx, query, blog.Title, blog.Content, blog.ID)
	return err
}

func (r *blogRepository) Delete(ctx context.Context, id int64) error {
	query := "UPDATE blog SET deleted_at = ? WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *blogRepository) List(ctx context.Context, limit, offset int) ([]*model.Blog, error) {
	query := "SELECT id, user_id, title, content " +
		"FROM blog " +
		"WHERE deleted_at IS NULL " +
		"ORDER BY id DESC " +
		"LIMIT ? OFFSET ?"

	rows, err := r.db.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []*model.Blog
	for rows.Next() {
		blog := &model.Blog{}
		if err := rows.Scan(&blog.ID, &blog.UserID, &blog.Title, &blog.Content); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}
