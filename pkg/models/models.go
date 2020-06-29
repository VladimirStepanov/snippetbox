package models

import (
	"errors"
	"time"
)

//Custom errors
var (
	ErrNoRecord = errors.New("Record not found")
)

//User model for users table
type User struct {
	ID             int
	Name           string
	Surname        string
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
