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

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("Error! %d != %d", w.Result().StatusCode, http.StatusInternalServerError)
	}
}

func TestRenderSuccess(t *testing.T) {
	s := NewTestServer()
	var err error
	s.templateCache, err = newTemplateCache("../../ui/html")

	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()

	s.render(w, "home")

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Error! %d != %d", w.Result().StatusCode, http.StatusOK)
	}
}
