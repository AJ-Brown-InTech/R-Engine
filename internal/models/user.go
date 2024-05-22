package models

import (
	"context"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type User struct {
    UserID              string    `json:"user_id" db:"user_id"`
    Username            string    `json:"username" db:"username"`
    UserPassword        string    `json:"-" db:"user_password"` 
    Email               string    `json:"email" db:"email"`
    EmailVerified       bool      `json:"email_verified" db:"email_verified"`
    FirstName           string    `json:"first_name,omitempty" db:"first_name"` 
    LastName            string    `json:"last_name,omitempty" db:"last_name"`   
    UserBio             string    `json:"user_bio,omitempty" db:"user_bio"`     
    Birthday            string    `json:"birthday,omitempty" db:"birthday"`     
    CreatedAt           time.Time `json:"created_at" db:"created_at"`
    UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
    Verified            bool      `json:"verified" db:"verified"`
    ProfilePictureURL   string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"` 
    NotificationsEnabled bool    `json:"notifications_enabled" db:"notifications_enabled"`
    Flagged             int       `json:"flagged" db:"flagged"`
    Rank                int       `json:"rank" db:"rank"`
    Creator             bool      `json:"creator" db:"creator"`
    Salt                []byte    `json:"-" db:"salt"`
    Latitude            float64   `json:"latitude" db:"latitude"`
    Longitude           float64   `json:"longitude" db:"longitude"`
    SessionToken        string    `json:"session_token" db:"session_token"`
}

// Check if username and email exist
func (u *User) UserAndEmailExist(ctx context.Context, db *sqlx.DB) (bool, bool, error) {
	var userCount, emailCount int

	// Check username
	err := db.GetContext(ctx, &userCount, "SELECT COUNT(*) FROM users WHERE username = $1", u.Username)
	if err != nil {
		return false, false, err
	}

	// Check email
	err = db.GetContext(ctx, &emailCount, "SELECT COUNT(*) FROM users WHERE email = $1", u.Email)
	if err != nil {
		return false, false, err
	}

	return userCount > 0, emailCount > 0, nil
}

// Check if username exist
func (u *User) UserExist(ctx context.Context, db *sqlx.DB) (bool, error) {
	var userCount int

	// Check username
	err := db.GetContext(ctx, &userCount, "SELECT COUNT(*) FROM users WHERE username = $1", u.Username)
	if err != nil {
		return false, err
	}

	return userCount > 0,  nil
}

// Check if email exist
func (u *User) EmailExist(ctx context.Context, db *sqlx.DB) (bool, error) {
	var emailCount int

	// Check email
	err := db.GetContext(ctx, &emailCount, "SELECT COUNT(*) FROM users WHERE email = $1", u.Email)
	if err != nil {
		return false, err
	}

	return emailCount > 0, nil
}

// Validate user creation
func (u *User) Validate() error {
	validation := validator.New()
    var UserValidation = map[string]string{
        "UserId":      "required,uuid",
        "Username":    "required,min=8,max=15",
        "Password":    "required,min=12",
        "Email":       "required,email",
        "FirstName":   "max=50",
        "LastName":    "max=50",
    }

	validation.RegisterStructValidationMapRules(UserValidation, User{})
    return  validation.Struct(&u)
}


type Following struct {
    FollowerID  string `json:"follower_id" db:"follower_id"`
    FollowingID string `json:"following_id" db:"following_id"`
}

