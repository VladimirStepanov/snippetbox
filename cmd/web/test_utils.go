package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func getTestUserData() map[int64]*models.User {
	um := map[int64]*models.User{}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), 14)
	um[1] = &models.User{ID: 1, Firstname: "Ivan", Lastname: "Doe", Email: "vova@mail.com", HashedPassword: hashedPassword}
	um[2] = &models.User{ID: 2, Firstname: "Conor", Lastname: "Ivanov", Email: "conor@mail.com", HashedPassword: hashedPassword}

	return um
}

func getTestSnippetData(startID, count int, isPub bool, oID int64) []*models.Snippet {
	ss := []*models.Snippet{}
	for i := startID; i < startID+count; i++ {
		ss = append(ss, &models.Snippet{
			ID:       int64(i),
			Title:    fmt.Sprintf("%dtitle%d", i, i),
			Content:  fmt.Sprintf("%content%d", i, i),
			Created:  time.Now(),
			Expires:  time.Now().Add(time.Hour),
			OwnerID:  oID,
			IsPublic: isPub,
		})
	}

	return ss
}

//NewTestServer return *Server test object
func NewTestServer(sr models.SnippetRepository, ur models.UserRepository) *Server {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	return New(":8080", logger, ur, sr, nil)
}

//NewTestServerWithUI return *Server object with templateCache
func NewTestServerWithUI(dir string, sr models.SnippetRepository, ur models.UserRepository) (*Server, error) {
	s := NewTestServer(sr, ur)
	var err error
	s.templateCache, err = newTemplateCache(dir)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func get(url string, t *testing.T, srv *httptest.Server) (int, http.Header, []byte) {
	rs, err := srv.Client().Get(url)

	if err != nil {
		t.Fatal(err.Error())
	}

	defer rs.Body.Close()

	data, err := ioutil.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, data
}
