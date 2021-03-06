package models

import (
	"errors"
	"time"
)

//Custom errors
var (
	ErrNoRecord       = errors.New("models: Record not found")
	ErrDuplicateEmail = errors.New("models: Duplicate email")
	ErrAuth           = errors.New("models: Can't find user in database")
	ErrUnknownOwnerID = errors.New("models: Unknown snippet owner ID ")
)

//User model for users table
type User struct {
	ID             int64
	Firstname      string
	Lastname       string
	Email          string
	Password       string
	HashedPassword []byte
	LogoutHash     string
}

//Snippet model for snippets table
type Snippet struct {
	ID       int64
	Title    string
	Content  string
	Created  time.Time
	Expires  time.Time
	OwnerID  int64
	IsPublic bool
}
