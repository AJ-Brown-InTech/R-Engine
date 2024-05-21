package models

import "time"


type  Session struct {
	Expiration *time.Time `json:"expiration" db:"expiration"`
	Username            string    `json:"username" db:"username"`
}


func (s Session) Create(user string) (*Session){
	// Add 48hour expiration
	var duration time.Duration = time.Hour * 48 
	ex :=time.Now().Add(duration)
	s.Expiration = &ex
	s.Username = user
	return &s
}


func (s *Session) Validate() (bool, error) {
	// TODO: 
	return false, nil
}


// TODO: delete session
// TODO: middlware to validate session