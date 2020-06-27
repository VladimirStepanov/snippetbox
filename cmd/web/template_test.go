package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRenderNotFound(t *testing.T) {
	s := NewTestServer()
	s.templateCache = map[string]*template.Template{}

	w := httptest.NewRecorder()

	s.render(w, "broken", &templateData{Title: "TestPage"})

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

	s.render(w, "snippets", &templateData{Title: "TestPage"})

	if w.Code != http.StatusOK {
		t.Fatalf("Error! %d != %d", w.Code, http.StatusOK)
	}
}

func TestAddDefaultTemplateData(t *testing.T) {
	title := "FFFFFPAGE"
	s, err := NewTestServerWithUI("../../ui/html")

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	s.render(w, "snippets", &templateData{Title: title})

	if !strings.Contains(w.Body.String(), title) {
		t.Fatalf("Body doesn't contain %s", title)
	}

	if !strings.Contains(w.Body.String(), fmt.Sprintf("%d", time.Now().Year())) {
		t.Fatalf("Body doesn't contain current year %d", time.Now().Year())
	}
}
