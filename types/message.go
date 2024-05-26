package types

import (
	"context"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

// Message represents a message between users
type Message struct {
	MessageID   string    `json:"message_id" db:"message_id"`
	SenderID    string    `json:"sender_id" db:"sender_id"`
	ReceiverID  string    `json:"receiver_id" db:"receiver_id"`
	ContentType string    `json:"content_type" db:"content_type"`
	Content     string    `json:"content,omitempty" db:"content"`
	MediaURL    string    `json:"media_url,omitempty" db:"media_url"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	IsRead      bool      `json:"is_read" db:"is_read"`
}

// Create a new message
func (m *Message) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO messages (message_id, sender_id, receiver_id, content_type, content, media_url, timestamp, is_read) VALUES (:message_id, :sender_id, :receiver_id, :content_type, :content, :media_url, :timestamp, :is_read)`
	_, err := db.NamedExecContext(ctx, query, m)
	return err
}

// Read a message by ID
func (m *Message) Read(ctx context.Context, db *sqlx.DB, messageID uuid.UUID) error {
	query := `SELECT * FROM messages WHERE message_id = $1`
	return db.GetContext(ctx, m, query, messageID)
}

// Update a message
func (m *Message) Update(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE messages SET content_type=:content_type, content=:content, media_url=:media_url, timestamp=:timestamp, is_read=:is_read WHERE message_id=:message_id`
	_, err := db.NamedExecContext(ctx, query, m)
	return err
}

// Delete a message by ID
func (m *Message) Delete(ctx context.Context, db *sqlx.DB, messageID uuid.UUID) error {
	query := `DELETE FROM messages WHERE message_id = $1`
	_, err := db.ExecContext(ctx, query, messageID)
	return err
}

// List messages between two users
func ListMessages(ctx context.Context, db *sqlx.DB, senderID, receiverID string) ([]Message, error) {
	var messages []Message
	query := `SELECT * FROM messages WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) ORDER BY timestamp DESC`
	err := db.SelectContext(ctx, &messages, query, senderID, receiverID)
	return messages, err
}
