package types

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// BlockedUser represents a blocked user relationship
type BlockedUser struct {
	BlockerID     string `json:"blocker_id" db:"blocker_id"`
	BlockedUserID string `json:"blocked_user_id" db:"blocked_user_id"`
}

// Create a new blocked user relationship
func (b *BlockedUser) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO blocked_users (blocker_id, blocked_user_id) VALUES (:blocker_id, :blocked_user_id)`
	_, err := db.NamedExecContext(ctx, query, b)
	return err
}

// Delete a blocked user relationship
func (b *BlockedUser) Delete(ctx context.Context, db *sqlx.DB) error {
	query := `DELETE FROM blocked_users WHERE blocker_id = $1 AND blocked_user_id = $2`
	_, err := db.ExecContext(ctx, query, b.BlockerID, b.BlockedUserID)
	return err
}

// List blocked users for a user
func ListBlockedUsers(ctx context.Context, db *sqlx.DB, userID string) ([]BlockedUser, error) {
	var blockedUsers []BlockedUser
	query := `SELECT * FROM blocked_users WHERE blocker_id = $1`
	err := db.SelectContext(ctx, &blockedUsers, query, userID)
	return blockedUsers, err
}
