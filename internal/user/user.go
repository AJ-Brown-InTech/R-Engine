package services

import (
	"Engine/internal/models"
	"Engine/internal/database"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/scrypt"
)

// Create a new account
func Create(ctx context.Context, db *sqlx.DB, u models.User) error {

// Check if username and email already exist
userExists, emailExists, err := u.UserAndEmailExist(ctx, db)
if err != nil {
return err
}

if userExists {
return fmt.Errorf("username %s already exists", u.Username)
}

if emailExists {
return fmt.Errorf("email %s already exists", u.Email)
}

// Create a new user account hash
hash, salt, err := Hash(u.UserPassword)
if err != nil {
return err
}

// setup user account default settings
u.UserID = uuid.New().String()
u.Verified = false
u.UserPassword = hash
u.Salt = salt
u.UserBio = ""


// validate 
err = u.Validate()
if err != nil {
return err
}

tx, err := db.Beginx()
if err != nil {
	return err
}
return database.InsertTransaction(ctx,tx,"users", u)	
}

// TODO: Login/Logout of account

// TODO: Delete account

// TODO: Update user account

// TODO: Follow/Unfollow user acccount

// TODO: Block/unblock user account

// TODO: Fetch user account

// TODO: flag account

// TODO: suspend account




func Hash(password string) (string, []byte, error) {
	salt := []byte(password)
	hash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return "", nil, err
	}
	return base64.StdEncoding.EncodeToString(hash), salt, nil
}

// Check if the provided password matches the stored hashed password
func Compare(hash string, salt []byte, input string) bool {
	ihash, err := scrypt.Key([]byte(input), salt, 16384, 8, 1, 32)
	if err != nil {
		return false
	}
	eihash := base64.StdEncoding.EncodeToString(ihash)
	return hash == eihash
}
