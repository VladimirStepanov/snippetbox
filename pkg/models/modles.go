package models

import "time"

//User model for users table
type User struct {
	ID       int
	Name     string
	Surname  string
	Email    string
	Password []byte
}

//Snippet model for snippets table
type Snippet struct {
	ID       int
	Title    string
	Content  string
	Created  time.Time
	Expires  time.Time
	OwnerID  int
	IsPublic bool
}
