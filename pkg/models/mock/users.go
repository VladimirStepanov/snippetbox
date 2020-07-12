package mock

import (
	"math/rand"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

//UsersStore mock for test endpoint
type UsersStore struct {
	DB map[int64]*models.User
}

func getRandUserID(m map[int64]*models.User) int64 {

	for {
		id := rand.Int63()
		if _, ok := m[id]; !ok {
			return id
		}
	}
}

//Insert user in map
func (us *UsersStore) Insert(firstname, lastname, mail, password string) (int64, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return 0, err
	}

	for _, value := range us.DB {
		if value.Email == mail {
			return 0, models.ErrDuplicateEmail
		}
	}

	id := getRandUserID(us.DB)

	us.DB[id] = &models.User{ID: id, Firstname: firstname, Lastname: lastname, Email: mail, HashedPassword: hashedPassword}

	return id, nil
}

//Get User from map
func (us *UsersStore) Get(id int64) (*models.User, error) {
	if val, ok := us.DB[id]; ok {
		retVal := &models.User{}
		*retVal = *val
		retVal.HashedPassword = nil
		return retVal, nil
	}

	return nil, models.ErrNoRecord
}

//Authenticate ...
func (us *UsersStore) Authenticate(email, password string) (int64, error) {
	for id, value := range us.DB {
		if value.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(value.HashedPassword), []byte(password))
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return 0, models.ErrAuth
			} else if err != nil {
				return 0, err
			}
			return id, nil
		}
	}

	return 0, models.ErrAuth

}
