package types

import (
	"context"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

// Comment represents a comment on a post
type Comment struct {
	CommentID string    `json:"comment_id" db:"comment_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	PostID    string    `json:"post_id" db:"post_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Create a new comment
func (c *Comment) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO comments (comment_id, user_id, post_id, content, created_at) VALUES (:comment_id, :user_id, :post_id, :content, :created_at)`
	_, err := db.NamedExecContext(ctx, query, c)
	return err
}

// Read a comment by ID
func (c *Comment) Read(ctx context.Context, db *sqlx.DB, commentID uuid.UUID) error {
	query := `SELECT * FROM comments WHERE comment_id = $1`
	return db.GetContext(ctx, c, query, commentID)
}

// Update a comment
func (c *Comment) Update(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE comments SET content=:content, created_at=:created_at WHERE comment_id=:comment_id`
	_, err := db.NamedExecContext(ctx, query, c)
	return err
}

// Delete a comment by ID
func (c *Comment) Delete(ctx context.Context, db *sqlx.DB, commentID uuid.UUID) error {
	query := `DELETE FROM comments WHERE comment_id = $1`
	_, err := db.ExecContext(ctx, query, commentID)
	return err
}

// List comments for a post
func ListComments(ctx context.Context, db *sqlx.DB, postID string) ([]Comment, error) {
	var comments []Comment
	query := `SELECT * FROM comments WHERE post_id = $1 ORDER BY created_at DESC`
	err := db.SelectContext(ctx, &comments, query, postID)
	return comments, err
}
