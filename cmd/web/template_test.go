package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models/mock"
)

func TestRenderNotFound(t *testing.T) {
	s := NewTestServer(&mock.SnippetStore{}, &mock.UsersStore{})
	s.templateCache = map[string]*template.Template{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	s.render(w, r, "broken", &templateData{Title: "TestPage"})

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Error! %d != %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRenderSuccess(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{})

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	s.render(w, r, "snippets", &templateData{Title: "TestPage"})

	if w.Code != http.StatusOK {
		t.Fatalf("Error! %d != %d", w.Code, http.StatusOK)
	}
}

func TestAddDefaultTemplateData(t *testing.T) {
	title := "FFFFFPAGE"
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{})

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	s.render(w, r, "snippets", &templateData{Title: title})

	if !strings.Contains(w.Body.String(), title) {
		t.Fatalf("Body doesn't contain %s", title)
	}

	if !strings.Contains(w.Body.String(), fmt.Sprintf("%d", time.Now().Year())) {
		t.Fatalf("Body doesn't contain current year %d", time.Now().Year())
	}
}
