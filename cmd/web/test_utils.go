package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

//NewTestServer return *Server test object
func NewTestServer() *Server {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	return New(":8080", logger)
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
