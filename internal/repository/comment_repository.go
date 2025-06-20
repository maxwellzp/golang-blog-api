package repository

import (
	"context"
	"database/sql"
	"maxwellzp/blog-api/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id int64) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id int64) error
	ListByBlogID(ctx context.Context, blogID int64, limit, offset int) ([]*model.Comment, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := "INSERT INTO comment (user_id, blog_id, content) VALUES (?, ?, ?)"
	res, err := r.db.ExecContext(ctx, query, comment.UserID, comment.BlogID, comment.Content)
	if err != nil {
		return err
	}
	comment.ID, err = res.LastInsertId()
	return err
}

func (r *commentRepository) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	query := "SELECT id, user_id, blog_id, content FROM comment WHERE id = ?"

	row := r.db.QueryRowContext(ctx, query, id)

	comment := &model.Comment{}
	if err := row.Scan(&comment.ID, &comment.UserID, &comment.BlogID, &comment.Content); err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *model.Comment) error {
	query := "UPDATE comment SET content = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, comment.Content, comment.ID)
	return err
}

func (r *commentRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM comment WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *commentRepository) ListByBlogID(ctx context.Context, blogID int64, limit, offset int) ([]*model.Comment, error) {
	query := "SELECT id, user_id, blog_id, content " +
		"FROM comment " +
		"WHERE blog_id = ?" +
		" ORDER BY id DESC " +
		"LIMIT ? OFFSET ?"

	rows, err := r.db.QueryContext(ctx, query, blogID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(&c.ID, &c.UserID, &c.BlogID, &c.Content); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
