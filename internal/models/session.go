package models

import (
	"encoding/base64"
	"encoding/json"
	"time"
	"github.com/go-playground/validator/v10"
)

type Session struct {
	Token string
}

type Token struct {
	User       string
	Expiration time.Time
}

func (t *Token) Create(uid string) (*Token, error) {

	t.Expiration = time.Now().Add(72 * time.Hour) // We can change if needed
	t.User = uid

	// Validate the session struct
	var TokenValidation = map[string]string{
		"User":       "required,uuid",
		"Expiration": "required,datetime",
	}

	validation := validator.New()
	validation.RegisterStructValidationMapRules(TokenValidation, Token{})
	if err := validation.Struct(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s Session) Create(t *Token) (*Session, error) {
	// Validate the session struct
	var TokenValidation = map[string]string{
		"User":       "required,uuid",
		"Expiration": "required,datetime",
	}

	validation := validator.New()
	validation.RegisterStructValidationMapRules(TokenValidation, Token{})
	if err := validation.Struct(t); err != nil {
		return nil, err
	}

	// Convert session to JSON
	sessionObj, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	// Encode the JSON as base64
	sessionB64 := base64.StdEncoding.EncodeToString(sessionObj)

	session := Session{Token: sessionB64}

	return &session, nil
}

func (s *Session) Validate() (bool, error) {
	// Decode base64
	tokenJSON, err := base64.StdEncoding.DecodeString(s.Token)
	if err != nil {
		return false, err
	}

	var token Token
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return false, err
	}

	// Check if token expiration is before current time
	return token.Expiration.Before(time.Now()), nil
}

