package types

import (
	"context"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

// Notification represents a notification for a user
type Notification struct {
	NotificationID string    `json:"notification_id" db:"notification_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	Type           string    `json:"type" db:"type"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	IsRead         bool      `json:"is_read" db:"is_read"`
}

// Create a new notification
func (n *Notification) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO notifications (notification_id, user_id, type, content, created_at, is_read) VALUES (:notification_id, :user_id, :type, :content, :created_at, :is_read)`
	_, err := db.NamedExecContext(ctx, query, n)
	return err
}

// Read a notification by ID
func (n *Notification) Read(ctx context.Context, db *sqlx.DB, notificationID uuid.UUID) error {
	query := `SELECT * FROM notifications WHERE notification_id = $1`
	return db.GetContext(ctx, n, query, notificationID)
}

// Update a notification
func (n *Notification) Update(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE notifications SET type=:type, content=:content, created_at=:created_at, is_read=:is_read WHERE notification_id=:notification_id`
	_, err := db.NamedExecContext(ctx, query, n)
	return err
}

// Delete a notification by ID
func (n *Notification) Delete(ctx context.Context, db *sqlx.DB, notificationID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE notification_id = $1`
	_, err := db.ExecContext(ctx, query, notificationID)
	return err
}

// List notifications for a user
func ListNotifications(ctx context.Context, db *sqlx.DB, userID string) ([]Notification, error) {
	var notifications []Notification
	query := `SELECT * FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	err := db.SelectContext(ctx, &notifications, query, userID)
	return notifications, err
}

// Mark a notification as read
func (n *Notification) MarkAsRead(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE notification_id = $1`
	_, err := db.ExecContext(ctx, query, n.NotificationID)
	return err
}
