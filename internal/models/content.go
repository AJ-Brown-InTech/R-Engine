package models

import "time"

type Post struct {
    PostID       string    `json:"post_id" db:"post_id"`
    UserID       string    `json:"user_id" db:"user_id"`
    Content      string    `json:"content,omitempty" db:"content"`
    PhotoURL     string    `json:"photo_url,omitempty" db:"photo_url"`
    Caption      string    `json:"caption,omitempty" db:"caption"`
    Latitude     float64   `json:"latitude" db:"latitude"`
    Longitude    float64   `json:"longitude" db:"longitude"`
    LocationName string    `json:"location_name" db:"location_name"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Like struct {
    LikeID    string    `json:"like_id" db:"like_id"`
    UserID    string    `json:"user_id" db:"user_id"`
    PostID    string    `json:"post_id" db:"post_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Comment struct {
    CommentID string    `json:"comment_id" db:"comment_id"`
    UserID    string    `json:"user_id" db:"user_id"`
    PostID    string    `json:"post_id" db:"post_id"`
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PostTag struct {
    PostID string `json:"post_id" db:"post_id"`
    TagID  string `json:"tag_id" db:"tag_id"`
}