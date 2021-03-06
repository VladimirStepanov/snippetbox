package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="gorilla.csrf.Token" value="(.+)">`)
var logoutHashRX = regexp.MustCompile(`/user/logout\?hash=(.+)'`)

type showSnippetsData struct {
	WantCode int
	WantSee  []*models.Snippet
	WantHide []*models.Snippet
	WantPage string
}

type showSnippetData struct {
	WantCode    int
	WantID      int64
	WantSnippet *models.Snippet
}

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

func setClearCookieJar(t *testing.T, srv *httptest.Server) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.Client().Jar = jar
}

//NewHttptestServer return new *httptest.Server object with cookie jar
func NewHttptestServer(t *testing.T, routes http.Handler) *httptest.Server {
	srv := httptest.NewServer(routes)

	setClearCookieJar(t, srv)

	return srv
}

//NewTestServer return *Server test object
func NewTestServer(sr models.SnippetRepository, ur models.UserRepository) *Server {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	testConfig := &Config{addr: ":8080", log: logger, sessionStore: sessions.NewCookieStore([]byte("123")), csrfKey: "123"}
	return New(testConfig, ur, sr)
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

func postForm(formData url.Values, url string, t *testing.T, srv *httptest.Server) (int, http.Header, []byte) {

	client := srv.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	rs, err := client.PostForm(url, formData)

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

func extractLogoutHash(t *testing.T, body []byte) string {
	return extractFromRE(logoutHashRX, t, body)
}

func extractCSRFToken(t *testing.T, body []byte) string {
	return extractFromRE(csrfTokenRX, t, body)
}

func extractFromRE(compRE *regexp.Regexp, t *testing.T, body []byte) string {
	matches := compRE.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no RE found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

func login(t *testing.T, srv *httptest.Server, email, password string) {

	code, _, body := get(fmt.Sprintf("%s/user/login", srv.URL), t, srv)

	if code != http.StatusOK {
		t.Fatalf("Return code %d != %d for csrf", code, http.StatusOK)
	}

	csrfToken := extractCSRFToken(t, body)

	formValues := url.Values{}
	formValues.Add("email", email)
	formValues.Add("password", password)
	formValues.Add("gorilla.csrf.Token", csrfToken)

	code, _, body = postForm(formValues, fmt.Sprintf("%s/user/login", srv.URL), t, srv)

	if code != http.StatusSeeOther {
		t.Fatalf("Return code %d != %d for postForm", code, http.StatusSeeOther)
	}

}

func testSnippetsPage(srv *httptest.Server, t *testing.T, tests map[string]showSnippetsData, path string) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, _, data := get(fmt.Sprintf("%s%s?page=%s", srv.URL, path, test.WantPage), t, srv)

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

func testShowSnippetPage(t *testing.T, srv *httptest.Server, tests map[string]showSnippetData) {
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
