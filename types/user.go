package types

import (
	"context"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

type User struct {
    UserID               string    `json:"user_id" db:"user_id"`
    Username             string    `json:"username" db:"username"`
    UserPassword         string    `json:"-" db:"user_password"` 
    Email                string    `json:"email" db:"email"`
    EmailVerified        bool      `json:"email_verified" db:"email_verified"`
    FirstName            string    `json:"first_name,omitempty" db:"first_name"` 
    LastName             string    `json:"last_name,omitempty" db:"last_name"`   
    UserBio              string    `json:"user_bio,omitempty" db:"user_bio"`     
    Birthday             string    `json:"birthday,omitempty" db:"birthday"`     
    CreatedAt            time.Time `json:"created_at" db:"created_at"`
    UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
    Verified             bool      `json:"verified" db:"verified"`
    ProfilePictureURL    string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"` 
    NotificationsEnabled bool      `json:"notifications_enabled" db:"notifications_enabled"`
    Flagged              int       `json:"flagged" db:"flagged"`
    Rank                 int       `json:"rank" db:"rank"`
    Creator              bool      `json:"creator" db:"creator"`
    Salt                 []byte    `json:"-" db:"salt"`
    Latitude             float64   `json:"latitude" db:"latitude"`
    Longitude            float64   `json:"longitude" db:"longitude"`
    SessionToken         string    `json:"session_token" db:"session_token"`
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
    return  validation.Struct(u)
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

// Check if username exists
func (u *User) UserExist(ctx context.Context, db *sqlx.DB) (bool, error) {
	var userCount int

	// Check username
	err := db.GetContext(ctx, &userCount, "SELECT COUNT(*) FROM users WHERE username = $1", u.Username)
	if err != nil {
		return false, err
	}

	return userCount > 0, nil
}

// Check if email exists
func (u *User) EmailExist(ctx context.Context, db *sqlx.DB) (bool, error) {
	var emailCount int

	// Check email
	err := db.GetContext(ctx, &emailCount, "SELECT COUNT(*) FROM users WHERE email = $1", u.Email)
	if err != nil {
		return false, err
	}

	return emailCount > 0, nil
}

// Create a new user
func (u *User) Create(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO users (user_id, username, user_password, email, email_verified, first_name, last_name, user_bio, birthday, created_at, updated_at, verified, profile_picture_url, notifications_enabled, flagged, rank, creator, salt, latitude, longitude, session_token)
			  VALUES (:user_id, :username, :user_password, :email, :email_verified, :first_name, :last_name, :user_bio, :birthday, :created_at, :updated_at, :verified, :profile_picture_url, :notifications_enabled, :flagged, :rank, :creator, :salt, :latitude, :longitude, :session_token)`
	
	_, err := db.NamedExecContext(ctx, query, u)
	return err
}

// Read a user by ID
func (u *User) Read(ctx context.Context, db *sqlx.DB, userID uuid.UUID) error {
	query := `SELECT * FROM users WHERE user_id = $1`
	return db.GetContext(ctx, u, query, userID)
}

// Update a user
func (u *User) Update(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE users SET username=:username, user_password=:user_password, email=:email, email_verified=:email_verified, first_name=:first_name, last_name=:last_name, user_bio=:user_bio, birthday=:birthday, updated_at=:updated_at, verified=:verified, profile_picture_url=:profile_picture_url, notifications_enabled=:notifications_enabled, flagged=:flagged, rank=:rank, creator=:creator, salt=:salt, latitude=:latitude, longitude=:longitude, session_token=:session_token
			  WHERE user_id=:user_id`
	
	_, err := db.NamedExecContext(ctx, query, u)
	return err
}

// Delete a user by ID
func (u *User) Delete(ctx context.Context, db *sqlx.DB, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := db.ExecContext(ctx, query, userID)
	return err
}

// List all users with pagination
func ListUsers(ctx context.Context, db *sqlx.DB, limit, offset int) ([]User, error) {
	var users []User
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}

// Search users by username or email
func SearchUsers(ctx context.Context, db *sqlx.DB, searchTerm string) ([]User, error) {
	var users []User
	query := `SELECT * FROM users WHERE username ILIKE $1 OR email ILIKE $1`
	err := db.SelectContext(ctx, &users, query, "%"+searchTerm+"%")
	return users, err
}
