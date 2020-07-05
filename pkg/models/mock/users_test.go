package mock

import (
	"reflect"
	"testing"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
)

func TestInsertUser(t *testing.T) {
	us := UsersStore{DB: map[int64]*models.User{}}

	_, err := us.Insert("test", "test", "test", "test")

	if err != nil {
		t.Fatal(err)
	}

}

func TestDuplicateEmail(t *testing.T) {
	us := UsersStore{DB: map[int64]*models.User{}}

	_, err := us.Insert("test", "test", "test", "test")

	if err != nil {
		t.Fatal(err)
	}

	_, err = us.Insert("test", "test", "test", "test")

	if err != models.ErrDuplicateEmail {
		t.Fatalf("get: %v, want: %v", err, models.ErrDuplicateEmail)
	}
}

func TestGetUser(t *testing.T) {
	tests := map[string]struct {
		UserID    int64
		WantUser  *models.User
		WantError error
	}{
		"No user": {
			UserID:    0,
			WantUser:  nil,
			WantError: models.ErrNoRecord,
		},
		"Get user success": {
			WantUser: &models.User{
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john@gmail.com",
			},
			WantError: nil,
		},
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			us := UsersStore{DB: map[int64]*models.User{}}
			if value.WantUser != nil {
				id, err := us.Insert(value.WantUser.Firstname, value.WantUser.Lastname, value.WantUser.Email, "1234")
				if err != nil {
					t.Fatal(err)
				}
				value.WantUser.ID = id
				value.UserID = id
			}

			resUser, err := us.Get(value.UserID)

			if value.WantError != nil && value.WantError != err {
				t.Fatalf("Error %v != %v", value.WantError, err)
			}

			if !reflect.DeepEqual(resUser, value.WantUser) {
				t.Fatalf("Want: %v, get: %v, err: %v", resUser, value.WantUser, err)
			}

		})
	}
}

func TestAuthentication(t *testing.T) {

	type TestData struct {
		Email        string
		Password     string
		WantPassword string
		WantAdd      bool
	}

	tests := map[string]struct {
		Data      *TestData
		WantError error
		WantID    int64
	}{
		"Email not found": {
			Data:      &TestData{WantAdd: false, Email: "hello@test.com"},
			WantError: models.ErrAuth,
		},
		"Bad password": {
			Data:      &TestData{WantAdd: true, Email: "hello@mail.com", Password: "123", WantPassword: "1234"},
			WantError: models.ErrAuth,
		},
		"Auth success": {
			Data:      &TestData{WantAdd: true, Email: "hello@mail.com", Password: "123", WantPassword: "123"},
			WantError: nil,
		},
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			us := UsersStore{DB: map[int64]*models.User{}}

			if value.Data.WantAdd {
				var err error
				value.WantID, err = us.Insert("1", "2", value.Data.Email, value.Data.Password)
				if err != nil {
					t.Fatal(err)
				}
			}

			authID, err := us.Authenticate(value.Data.Email, value.Data.WantPassword)

			if value.WantError != nil && value.WantError != err {
				t.Fatalf("want: %v, get: %v", value.WantError, err)
			}

			if value.WantError == nil {
				if value.WantID != authID {
					t.Fatalf("Want id: %d, Get id: %d", value.WantID, authID)
				}
			}
		})
	}
}
