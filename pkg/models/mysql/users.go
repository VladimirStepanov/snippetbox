package mysql

import (
	"database/sql"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//UsersStore struct for working with snippets table
type UsersStore struct {
	DB *sql.DB
}

//Insert user to database
func (us *UsersStore) Insert(firstname, lastname, mail, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return 0, err
	}

	res, err := us.DB.Exec(
		"INSERT INTO users (firstname, lastname, mail, password) VALUES (?, ?, ?, ?)",
		firstname,
		lastname,
		mail,
		hashedPassword,
	)

	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			if me.Number == 1062 {
				return 0, models.ErrDuplicateEmail
			}
		}
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

//Get user from database
func (us *UsersStore) Get(id int64) (*models.User, error) {
	resUser := &models.User{}
	row := us.DB.QueryRow("SELECT id, firstname, lastname, mail FROM users where id = ?", id)

	err := row.Scan(&resUser.ID, &resUser.Firstname, &resUser.Lastname, &resUser.Email)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return resUser, nil
}

//Authenticate ...
func (us *UsersStore) Authenticate(email, password string) (int64, error) {
	var returnID int64
	var hashedPassword string

	row := us.DB.QueryRow(`select id, password from users where mail=?`, email)

	err := row.Scan(&returnID, &hashedPassword)

	if err == sql.ErrNoRows {
		return 0, models.ErrAuth
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrAuth
	} else if err != nil {
		return 0, err
	}

	return returnID, nil

}
