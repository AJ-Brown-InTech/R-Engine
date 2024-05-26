package types

import (
	"context"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

// Like represents a like on a post
type Like struct {
	LikeID    string    `json:"like_id" db:"like_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	PostID    string    `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Create a new like
func (l *Like) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO likes (like_id, user_id, post_id, created_at) VALUES (:like_id, :user_id, :post_id, :created_at)`
	_, err := db.NamedExecContext(ctx, query, l)
	return err
}

// Read a like by ID
func (l *Like) Read(ctx context.Context, db *sqlx.DB, likeID uuid.UUID) error {
	query := `SELECT * FROM likes WHERE like_id = $1`
	return db.GetContext(ctx, l, query, likeID)
}

// Delete a like by ID
func (l *Like) Delete(ctx context.Context, db *sqlx.DB, likeID uuid.UUID) error {
	query := `DELETE FROM likes WHERE like_id = $1`
	_, err := db.ExecContext(ctx, query, likeID)
	return err
}

// List likes for a post
func ListLikesForPost(ctx context.Context, db *sqlx.DB, postID string) ([]Like, error) {
	var likes []Like
	query := `SELECT * FROM likes WHERE post_id = $1 ORDER BY created_at DESC`
	err := db.SelectContext(ctx, &likes, query, postID)
	return likes, err
}

// List likes by a user
func ListLikesByUser(ctx context.Context, db *sqlx.DB, userID string) ([]Like, error) {
	var likes []Like
	query := `SELECT * FROM likes WHERE user_id = $1 ORDER BY created_at DESC`
	err := db.SelectContext(ctx, &likes, query, userID)
	return likes, err
}
