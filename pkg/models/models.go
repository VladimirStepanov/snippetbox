package models

import (
	"errors"
	"time"
)

//Custom errors
var (
	ErrNoRecord       = errors.New("models: Record not found")
	ErrDuplicateEmail = errors.New("models: Duplicate email")
)

//User model for users table
type User struct {
	ID             int64
	Firstname      string
	Lastname       string
	Email          string
	Password       string
	HashedPassword []byte
}

//Snippet model for snippets table
type Snippet struct {
	ID       string
	Title    string
	Content  string
	Created  time.Time
	Expires  time.Time
	OwnerID  int
	IsPublic bool
}
