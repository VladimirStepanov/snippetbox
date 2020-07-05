package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"githib.com/VladimirStepanov/snippetbox/pkg/models/mock"
	"github.com/sirupsen/logrus"
)

//NewTestServer return *Server test object
func NewTestServer() *Server {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	us := map[int64]*models.User{}
	return New(":8080", logger, &mock.UsersStore{DB: us}, &mock.SnippetStore{DB: []*models.Snippet{}, UsersMap: us})
}

//NewTestServerWithUI return *Server object with templateCache
func NewTestServerWithUI(dir string) (*Server, error) {
	s := NewTestServer()
	var err error
	s.templateCache, err = newTemplateCache(dir)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func get(t *testing.T, srv *httptest.Server) (int, http.Header, []byte) {
	rs, err := srv.Client().Get(srv.URL)

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
