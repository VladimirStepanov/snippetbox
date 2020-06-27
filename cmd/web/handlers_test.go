package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	s := NewTestServer()
	var err error
	s.templateCache, err = newTemplateCache("../../ui/html")

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	code, _, _ := getTest(t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d", code, http.StatusOK)
	}
}
