package models

//UserRepository interface for working with DB
type UserRepository interface {
	Insert(firstname, lastname, mail, password string) (int64, error)
	Get(id int64) (*User, error)
	Authenticate(email, password string) (int64, error)
}

//SnippetRepository interface for working with DB
type SnippetRepository interface {
	Insert(title, content string, expire int, isPublic bool, ownerID int64) (int64, error)
	Delete(snippetID, userID int64) error
	Get(snippetID int64) (*Snippet, error)
	Update(snippet *Snippet, ownerID int64) error
	LatestAll(ownerID int64, count, page int) ([]*Snippet, error)
}
