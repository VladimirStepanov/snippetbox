package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"githib.com/VladimirStepanov/snippetbox/pkg/models/mock"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
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

func TestLoginOK(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{})

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	code, _, _ := get(fmt.Sprintf("%s/user/login", srv.URL), t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d", code, http.StatusOK)
	}
}

// tests for endpoint /
func TestHomeWithData(t *testing.T) {

	um := getTestUserData()
	ss := getTestSnippetData(1, 15, true, 1)

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(s.routes())

	defer srv.Close()

	tests := map[string]showSnippetsData{
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

	testSnippetsPage(srv, t, tests, "/")
}

//tests for endpoint /snippets
func TestUserSnippetsWithData(t *testing.T) {

	um := getTestUserData()
	ss := append(getTestSnippetData(1, 5, false, 2), getTestSnippetData(6, 3, false, 1)...)

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	login(t, srv, "conor@mail.com", "12345678")

	tests := map[string]showSnippetsData{
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
			WantSee:  ss[:5],
			WantHide: ss[5:],
			WantPage: "1",
		},
	}
	testSnippetsPage(srv, t, tests, "/snippets")
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

	tests := map[string]showSnippetData{
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

	testShowSnippetPage(t, srv, tests)
}

func TestShowSnippetForAuthUser(t *testing.T) {
	um := getTestUserData()
	ss := append(getTestSnippetData(1, 1, false, 2), getTestSnippetData(3, 1, false, 1)...)

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	login(t, srv, "conor@mail.com", "12345678")

	tests := map[string]showSnippetData{
		"ShowPrivateSnippet": {
			WantCode:    http.StatusOK,
			WantID:      ss[0].ID,
			WantSnippet: ss[0],
		},
		"Private snippet other user": {
			WantCode: http.StatusNotFound,
			WantID:   ss[1].ID,
		},
	}

	testShowSnippetPage(t, srv, tests)
}

func TestSignUpForm(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{DB: getTestUserData()})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	code, _, data := get(fmt.Sprintf("%s/user/signup", srv.URL), t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d", code, http.StatusOK)
	}

	csrfToken := extractCSRFToken(t, data)

	tests := map[string]struct {
		firstname string
		lastname  string
		email     string
		password  string
		WantCode  int
		WantData  []byte
		csrfToken string
	}{
		"Empty firstname":           {"", "1", "vova@mail.com", "123", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Empty lastname":            {"1", "", "vova@mail.com", "123", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Empty email":               {"1", "2", "", "123", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Empty password":            {"1", "2", "v@mail.com", "", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Bad email":                 {"1", "2", "v@a", "123", http.StatusOK, []byte("must be a valid email address"), csrfToken},
		"Email already exists":      {"1", "2", "conor@mail.com", "123345678", http.StatusOK, []byte("email already exists"), csrfToken},
		"Show password":             {"1", "2", "v@a", "123", http.StatusOK, []byte("the length must be between 8 and 20"), csrfToken},
		"Bad csrf":                  {"Ivan", "1", "vova@mail.com", "123", http.StatusForbidden, nil, "bad"},
		"User successfully created": {"Ivan", "1", "vova23@mail.com", "12345678", http.StatusSeeOther, nil, csrfToken},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			formValues := url.Values{}
			formValues.Add("firstname", test.firstname)
			formValues.Add("lastname", test.lastname)
			formValues.Add("email", test.email)
			formValues.Add("password", test.password)
			formValues.Add("gorilla.csrf.Token", test.csrfToken)

			code, _, body := postForm(formValues, fmt.Sprintf("%s/user/signup", srv.URL), t, srv)

			if code != test.WantCode {
				t.Fatalf("Want: %d, Get: %d", test.WantCode, code)
			}

			if code == http.StatusOK {
				if !bytes.Contains(body, test.WantData) {
					t.Fatalf("%s not in result body", string(test.WantData))
				}
			}
		})
	}
}

func TestLoginForm(t *testing.T) {
	um := getTestUserData()
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	getCsrf := func(s string) string {
		return s
	}

	badCsrf := func(s string) string {
		return "bad"
	}

	tests := map[string]struct {
		email        string
		password     string
		WantCode     int
		WantData     []byte
		getCsrfToken func(string) string
	}{
		"User successfully auth": {"conor@mail.com", "12345678", http.StatusSeeOther, nil, getCsrf},
		"Empty email":            {"", "123", http.StatusOK, []byte("cannot be blank"), getCsrf},
		"Empty password":         {"v@mail.com", "", http.StatusOK, []byte("cannot be blank"), getCsrf},
		"Bad csrf":               {"vova@mail.com", "123", http.StatusForbidden, nil, badCsrf},
		"Bad password":           {"conor@mail.com", "123", http.StatusOK, []byte("Email or password incorrect"), getCsrf},
		"Bad email":              {"conor1@mail.com", "12345678", http.StatusOK, []byte("Email or password incorrect"), getCsrf},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			setClearCookieJar(t, srv)

			code, _, data := get(fmt.Sprintf("%s/user/login", srv.URL), t, srv)

			if code != http.StatusOK {
				t.Fatalf("Return code %d != %d", code, http.StatusOK)
			}

			csrfToken := extractCSRFToken(t, data)

			formValues := url.Values{}
			formValues.Add("email", test.email)
			formValues.Add("password", test.password)
			formValues.Add("gorilla.csrf.Token", test.getCsrfToken(csrfToken))

			code, _, body := postForm(formValues, fmt.Sprintf("%s/user/login", srv.URL), t, srv)

			if code != test.WantCode {
				t.Fatalf("Want: %d, Get: %d", test.WantCode, code)
			}

			if code == http.StatusOK {
				if !bytes.Contains(body, test.WantData) {
					t.Fatalf("%s not in result body", string(test.WantData))
				}
			}
		})
	}
}

func TestAuthUserMiddleware(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{DB: getTestUserData()})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	code, _, body := get(srv.URL, t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d for home page", code, http.StatusOK)
	}

	if bytes.Contains(body, []byte("Logout")) || bytes.Contains(body, []byte("My snippets")) {
		t.Fatal("'Logout' or 'My snippets' on home page")
	}

	login(t, srv, "conor@mail.com", "12345678")

	code, _, body = get(srv.URL, t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d for home page", code, http.StatusOK)
	}

	if !bytes.Contains(body, []byte("Logout")) || !bytes.Contains(body, []byte("My snippets")) {
		t.Fatal("'Logout' or 'My snippets' not on home page")
	}
}

func TestAccessOnlyNotAuthMiddleware(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{DB: getTestUserData()})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	login(t, srv, "conor@mail.com", "12345678")

	tests := map[string]struct {
		Path     string
		WantCode int
	}{
		"login": {
			Path:     "/user/login",
			WantCode: http.StatusSeeOther,
		},
		"signup": {
			Path:     "/user/signup",
			WantCode: http.StatusSeeOther,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, _, _ := get(fmt.Sprintf("%s%s", srv.URL, test.Path), t, srv)
			if code != test.WantCode {
				t.Fatalf("Want: %d, Get: %d", test.WantCode, code)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{}, &mock.UsersStore{DB: getTestUserData()})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()

	tests := map[string]struct {
		Path       func(h string) string
		WantLogout bool
	}{
		"logout success": {
			Path: func(h string) string {
				return "/user/logout?hash=" + h
			},
			WantLogout: false,
		},
		"empty hash": {
			Path: func(h string) string {
				return "/user/logout"
			},
			WantLogout: true,
		},
		"bad logout hash": {
			Path: func(h string) string {
				return "/user/logout?hash=badhash"
			},
			WantLogout: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			setClearCookieJar(t, srv)
			login(t, srv, "conor@mail.com", "12345678")
			code, _, body := get(srv.URL, t, srv)

			if code != http.StatusOK {
				t.Fatalf("Return code %d != %d for home page", code, http.StatusOK)
			}

			logoutHash := extractLogoutHash(t, body)

			code, _, _ = get(fmt.Sprintf("%s%s", srv.URL, test.Path(logoutHash)), t, srv)
			if code != http.StatusSeeOther {
				t.Fatalf("Want: %d, Get: %d", http.StatusSeeOther, code)
			}

			code, _, body = get(srv.URL, t, srv)

			if code != http.StatusOK {
				t.Fatalf("Return code %d != %d for home page in tests", code, http.StatusOK)
			}

			if bytes.Contains(body, []byte("Logout")) != test.WantLogout {
				t.Fatalf("I want logout: %v, but not", test.WantLogout)
			}
		})
	}
}

func TestDelete(t *testing.T) {

	um := getTestUserData()
	ss := append(getTestSnippetData(1, 5, false, 2), getTestSnippetData(6, 3, false, 1)...)

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()
	login(t, srv, "conor@mail.com", "12345678")

	code, _, body := get(srv.URL, t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d for home page", code, http.StatusOK)
	}

	logoutHash := extractLogoutHash(t, body)

	tests := map[string]struct {
		WantID   int64
		WantCode int
		Path     func(h string, id int64) string
	}{
		"Delete with empty hash": {
			WantID:   ss[0].ID,
			WantCode: http.StatusNotFound,
			Path: func(h string, id int64) string {
				return fmt.Sprintf("/snippet/delete/%d", id)
			},
		},
		"Delete with bad hash": {
			WantID:   ss[0].ID,
			WantCode: http.StatusNotFound,
			Path: func(h string, id int64) string {
				return fmt.Sprintf("/snippet/delete/%d?hash=123", id)
			},
		},
		"Delete not my own snippet": {
			WantID:   ss[7].ID,
			WantCode: http.StatusNotFound,
			Path: func(h string, id int64) string {
				return fmt.Sprintf("/snippet/delete/%d?hash=%s", id, h)
			},
		},
		"Success delete": {
			WantID:   ss[0].ID,
			WantCode: http.StatusSeeOther,
			Path: func(h string, id int64) string {
				return fmt.Sprintf("/snippet/delete/%d?hash=%s", id, h)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, _, _ = get(fmt.Sprintf("%s%s", srv.URL, test.Path(logoutHash, test.WantID)), t, srv)

			if code != test.WantCode {
				t.Fatalf("Error! Want %d, get %d", test.WantCode, code)
			}
		})
	}
}

func TestServer_deleteSnippet(t *testing.T) {
	type fields struct {
		addr          string
		log           *logrus.Logger
		templateCache map[string]*template.Template
		userStore     models.UserRepository
		snippetStore  models.SnippetRepository
		session       *sessions.CookieStore
		csrfKey       string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				addr:          tt.fields.addr,
				log:           tt.fields.log,
				templateCache: tt.fields.templateCache,
				userStore:     tt.fields.userStore,
				snippetStore:  tt.fields.snippetStore,
				session:       tt.fields.session,
				csrfKey:       tt.fields.csrfKey,
			}
			s.deleteSnippet(tt.args.w, tt.args.r)
		})
	}
}

func TestCreateSnippetForm(t *testing.T) {
	um := getTestUserData()
	ss := append(getTestSnippetData(1, 5, false, 2))

	s, err := NewTestServerWithUI("../../ui/html", &mock.SnippetStore{DB: ss, UsersMap: um}, &mock.UsersStore{DB: um})

	if err != nil {
		t.Fatal(err)
	}

	srv := NewHttptestServer(t, s.routes())
	defer srv.Close()
	login(t, srv, "conor@mail.com", "12345678")

	code, _, data := get(fmt.Sprintf("%s/snippet/create", srv.URL), t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d", code, http.StatusOK)
	}

	csrfToken := extractCSRFToken(t, data)

	tests := map[string]struct {
		title       string
		content     string
		expire      string
		snippetType string
		WantCode    int
		WantData    []byte
		csrfToken   string
	}{
		"Bad csrf token":                    {"1", "2", "3", "4", http.StatusForbidden, nil, "bad"},
		"Empty title":                       {"", "2", "3", "4", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Empty content":                     {"1", "", "3", "4", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Empty expire":                      {"1", "2", "", "4", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Empty type":                        {"1", "2", "3", "", http.StatusOK, []byte("cannot be blank"), csrfToken},
		"Type not in ['Public', 'Private']": {"1", "2", "3", "Bad", http.StatusOK, []byte("must be a valid value"), csrfToken},
		"Negative expire":                   {"1", "2", "-3", "4", http.StatusOK, []byte("value must be integer greater than zero"), csrfToken},
		"Bad expire":                        {"1", "2", "ff", "4", http.StatusOK, []byte("value must be integer greater than zero"), csrfToken},
		"Success create private":            {"title", "content", "150", "Private", http.StatusSeeOther, nil, csrfToken},
		"Success create public":             {"title", "content", "150", "Public", http.StatusSeeOther, nil, csrfToken},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			formValues := url.Values{}
			formValues.Add("title", test.title)
			formValues.Add("content", test.content)
			formValues.Add("expire", test.expire)
			formValues.Add("type", test.snippetType)
			formValues.Add("gorilla.csrf.Token", test.csrfToken)

			code, _, body := postForm(formValues, fmt.Sprintf("%s/snippet/create", srv.URL), t, srv)

			if code != test.WantCode {
				t.Fatalf("Want: %d, Get: %d", test.WantCode, code)
			}

			if code == http.StatusOK {
				if !bytes.Contains(body, test.WantData) {
					t.Fatalf("%s not in result body", string(test.WantData))
				}
			}
		})
	}
}
