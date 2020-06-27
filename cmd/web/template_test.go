package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderNotFound(t *testing.T) {
	s := NewTestServer()
	s.templateCache = map[string]*template.Template{}

	w := httptest.NewRecorder()

	s.render(w, "broken")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Error! %d != %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRenderSuccess(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html")

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	s.render(w, "home")

	if w.Code != http.StatusOK {
		t.Fatalf("Error! %d != %d", w.Code, http.StatusOK)
	}
}
