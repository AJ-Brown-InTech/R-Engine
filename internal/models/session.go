package models

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

// Create a new session with 48-hour expiration and encode it as a base64 string
func (s *Session) Create(user string) (string, error) {
	duration := 48 * time.Hour
	ex := time.Now().Add(duration)
	s.Expiration = &ex
	s.Username = user

	sessionJSON, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	sessionB64 := base64.StdEncoding.EncodeToString(sessionJSON)
	return sessionB64, nil
}

// Validate the session token by decoding it and checking the expiration
func (s *Session) Validate(token string) (bool, error) {
	tokenJSON, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, err
	}

	var session Session
	if err := json.Unmarshal(tokenJSON, &session); err != nil {
		return false, err
	}

	if session.Expiration.Before(time.Now()) {
		return false, errors.New("session expired")
	}

	s.Expiration = session.Expiration
	s.Username = session.Username

	return true, nil
}
