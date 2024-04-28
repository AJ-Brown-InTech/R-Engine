package models

import (
    "time"
)

type Notification struct {
    NotificationID string    `json:"notification_id" db:"notification_id"`
    UserID         string    `json:"user_id" db:"user_id"`
    Type           string    `json:"type" db:"type"`
    Content        string    `json:"content" db:"content"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    IsRead         bool      `json:"is_read" db:"is_read"`
}
