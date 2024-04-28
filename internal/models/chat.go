package models

import "time"

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