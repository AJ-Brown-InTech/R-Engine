package types

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// FlaggedAccount represents a flagged account
type FlaggedAccount struct {
	UserID      string `json:"user_id" db:"user_id"`
	FlagCount   int    `json:"flag_count" db:"flag_count"`
	IsSuspended bool   `json:"is_suspended" db:"is_suspended"`
	Reason      string `json:"reason,omitempty" db:"reason"`
}

// Create a new flagged account
func (f *FlaggedAccount) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO flagged_accounts (user_id, flag_count, is_suspended, reason) VALUES (:user_id, :flag_count, :is_suspended, :reason)`
	_, err := db.NamedExecContext(ctx, query, f)
	return err
}

// Read a flagged account by user ID
func (f *FlaggedAccount) Read(ctx context.Context, db *sqlx.DB, userID string) error {
	query := `SELECT * FROM flagged_accounts WHERE user_id = $1`
	return db.GetContext(ctx, f, query, userID)
}

// Update a flagged account
func (f *FlaggedAccount) Update(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE flagged_accounts SET flag_count=:flag_count, is_suspended=:is_suspended, reason=:reason WHERE user_id=:user_id`
	_, err := db.NamedExecContext(ctx, query, f)
	return err
}

// Delete a flagged account by user ID
func (f *FlaggedAccount) Delete(ctx context.Context, db *sqlx.DB, userID string) error {
	query := `DELETE FROM flagged_accounts WHERE user_id = $1`
	_, err := db.ExecContext(ctx, query, userID)
	return err
}

// List all flagged accounts
func ListFlaggedAccounts(ctx context.Context, db *sqlx.DB) ([]FlaggedAccount, error) {
	var flaggedAccounts []FlaggedAccount
	query := `SELECT * FROM flagged_accounts ORDER BY flag_count DESC`
	err := db.SelectContext(ctx, &flaggedAccounts, query)
	return flaggedAccounts, err
}
