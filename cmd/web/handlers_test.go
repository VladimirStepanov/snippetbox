package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"githib.com/VladimirStepanov/snippetbox/pkg/models/mock"
)

func TestHomeHandler(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{})

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	code, _, data := get(srv.URL, t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d", code, http.StatusOK)
	}

	if !strings.Contains(string(data), "Snippets feed is empty") {
		t.Fatalf("No string 'Snippets feed is empty' in data %v\n", data)
	}
}

func TestSignUpOK(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{})

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	code, _, _ := get(fmt.Sprintf("%s/user/signup", srv.URL), t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d", code, http.StatusOK)
	}
}

func TestHomeHandlerWithData(t *testing.T) {

	um := getTestUserData()
	ss := getTestSnippetData(1, 15, true, 1)

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	tests := map[string]struct {
		WantCode int
		WantSee  []*models.Snippet
		WantHide []*models.Snippet
		WantPage string
	}{
		"Bad page number": {
			WantCode: 500,
			WantPage: "ff",
		},
		"Negative page number": {
			WantCode: 500,
			WantPage: "-1",
		},
		"Test first page": {
			WantCode: 200,
			WantSee:  ss[:10],
			WantHide: ss[10:],
			WantPage: "1",
		},
		"Test second page": {
			WantCode: 200,
			WantSee:  ss[11:],
			WantHide: ss[:10],
			WantPage: "2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, _, data := get(fmt.Sprintf("%s/?page=%s", srv.URL, test.WantPage), t, srv)

			if test.WantCode != code {
				t.Fatalf("Want code: %d, Get code: %d", test.WantCode, code)
			}

			if test.WantCode == http.StatusOK {
				for _, val := range test.WantSee {
					if !strings.Contains(string(data), val.Title) {
						t.Fatalf("Want see: %v", val)
					}
				}
				for _, val := range test.WantHide {
					if strings.Contains(string(data), val.Title) {
						t.Fatalf("Want hide:  %v", val)
					}
				}
			}
		})
	}
}

func TestShowSnippetForNotAuthUser(t *testing.T) {
	um := getTestUserData()
	// ss := getTestSnippetData(1, 15, true, 1)
	ss := append(getTestSnippetData(1, 1, true, 1), getTestSnippetData(5, 3, false, 1)...)

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	tests := map[string]struct {
		WantCode    int
		WantID      int64
		WantSnippet *models.Snippet
	}{
		"ShowSnippet": {
			WantCode:    http.StatusOK,
			WantID:      ss[0].ID,
			WantSnippet: ss[0],
		},
		"Snippet not found": {
			WantCode: http.StatusNotFound,
			WantID:   100500,
		},
		"Get private snippet": {
			WantCode: http.StatusNotFound,
			WantID:   ss[3].ID,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, _, data := get(fmt.Sprintf("%s/snippet/%d", srv.URL, test.WantID), t, srv)

			if test.WantCode != code {
				t.Fatalf("Want code: %d, Get code: %d", test.WantCode, code)
			}

			if test.WantCode == http.StatusOK {
				if !strings.Contains(string(data), test.WantSnippet.Title) || !strings.Contains(string(data), test.WantSnippet.Content) {
					t.Fatalf("Want see: %v", test.WantSnippet)
				}
			}
		})
	}
}
