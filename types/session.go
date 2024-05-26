package types

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Session struct {
	Expiration *time.Time `json:"expiration" db:"expiration"`
	Username   string     `json:"username" db:"username"`
}

// Create generates a new session with a 48-hour expiration and encodes it as a base64 string.
func (s *Session) Create(user string) (string, error) {
	// Set session duration to 48 hours
	duration := 48 * time.Hour
	expiration := time.Now().Add(duration)
	s.Expiration = &expiration
	s.Username = user

	// Marshal the session struct to JSON
	sessionJSON, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	// Encode the JSON as a base64 string
	sessionB64 := base64.StdEncoding.EncodeToString(sessionJSON)
	return sessionB64, nil
}

// Validate decodes the session token and checks if it has expired.
func (s *Session) Validate(token string) (bool, error) {
	// Decode the base64 encoded session token
	tokenJSON, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, errors.New("failed to decode session token")
	}

	// Unmarshal the JSON back to a Session struct
	var session Session
	if err := json.Unmarshal(tokenJSON, &session); err != nil {
		return false, errors.New("failed to unmarshal session JSON")
	}

	// Check if the session has expired
	if session.Expiration.Before(time.Now()) {
		return false, errors.New("session expired")
	}

	// Update the session fields
	s.Expiration = session.Expiration
	s.Username = session.Username

	return true, nil
}