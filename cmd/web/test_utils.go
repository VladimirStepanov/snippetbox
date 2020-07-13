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
	"testing"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="gorilla.csrf.Token" value="(.+)">`)

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
	return New(":8080", logger, ur, sr, sessions.NewCookieStore([]byte("123")), "123")
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

func extractCSRFToken(t *testing.T, body []byte) string {
	// Use the FindSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the
	// first position, and the values of any captured data in the subsequent
	// positions.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
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
