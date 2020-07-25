package main

import (
	"fmt"
	"net/http"
	"strconv"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func getPage(r *http.Request) (int, error) {
	page := 1
	pageStr := r.URL.Query().Get("page")

	if pageStr != "" {
		retPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return 0, err
		}

		if retPage < 1 {
			return 0, fmt.Errorf("Page less than 1")
		}

		page = retPage
	}

	return page, nil
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	page, err := getPage(r)
	if err != nil {
		s.serverError(w, err)
		return
	}
	snippets, err := s.snippetStore.LatestAll(-1, 10, page)

	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, r, "snippets", &templateData{Title: "Home", Snippets: snippets})
}

func (s *Server) userSnippets(w http.ResponseWriter, r *http.Request) {
	page, err := getPage(r)
	if err != nil {
		s.serverError(w, err)
		return
	}

	u := getAuthUserFromRequest(r)

	if u == nil {
		http.Redirect(w, r, "/", 303)
		return
	}

	snippets, err := s.snippetStore.LatestAll(u.ID, 10, page)

	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, r, "snippets", &templateData{Title: "My snippets", Snippets: snippets})
}

func (s *Server) showSnippet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	snippet, err := s.snippetStore.Get(int64(id))

	if err != nil {
		if err == models.ErrNoRecord {
			http.NotFound(w, r)
		} else {
			s.serverError(w, err)
		}
		return
	}

	currentUser := getAuthUserFromRequest(r)

	if !snippet.IsPublic {
		if currentUser == nil || currentUser.ID != snippet.OwnerID {
			http.NotFound(w, r)
			return
		}
	}

	var templateUser *models.User

	if currentUser != nil && snippet.OwnerID == currentUser.ID {
		templateUser = currentUser
	}

	s.render(w, r, "snippet", &templateData{Snippet: snippet, FormUser: templateUser})
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "signup", &templateData{CSRFField: csrf.TemplateField(r)})
}

func (s *Server) showLogin(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "login", &templateData{CSRFField: csrf.TemplateField(r)})
}

func (s *Server) signUpPOST(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
	}
	errors := validation.ValidateStruct(u,
		validation.Field(&u.Firstname, validation.Required),
		validation.Field(&u.Lastname, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 20)),
	)

	if errors != nil {
		errMap := errors.(validation.Errors)
		s.render(w, r, "signup", &templateData{Errors: errMap, FormUser: u, CSRFField: csrf.TemplateField(r)})
		return
	}
	_, err := s.userStore.Insert(u.Firstname, u.Lastname, u.Email, u.Password)

	if err == models.ErrDuplicateEmail {
		s.render(
			w, r,
			"signup",
			&templateData{
				Errors:    validation.Errors{"Email": fmt.Errorf("email already exists")},
				FormUser:  u,
				CSRFField: csrf.TemplateField(r),
			},
		)
		return
	} else if err != nil {
		s.serverError(w, err)
		return
	}

	if err := s.addFlashMessage(w, r, "User successfully created! Please log in.. "); err != nil {
		s.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/user/login", 303)

}

func (s *Server) loginPOST(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	errors := validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)

	if errors != nil {
		errMap := errors.(validation.Errors)
		s.render(w, r, "login", &templateData{Errors: errMap, FormUser: u, CSRFField: csrf.TemplateField(r)})
		return
	}

	userID, err := s.userStore.Authenticate(u.Email, u.Password)

	if err == models.ErrAuth {
		s.render(
			w, r,
			"login",
			&templateData{
				Errors:    validation.Errors{"Generic": fmt.Errorf("Email or password incorrect")},
				FormUser:  u,
				CSRFField: csrf.TemplateField(r),
			},
		)
		return
	} else if err != nil {
		s.serverError(w, err)
		return
	}

	if err = s.addNewUserSession(w, r, userID); err != nil {
		s.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", 303)

}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("hash")
	currentUser := getAuthUserFromRequest(r)

	if currentUser.LogoutHash == hash {
		session, err := s.session.Get(r, "SID")
		if err != nil {
			s.serverError(w, err)
			return
		}
		removeSession(w, r, session)
		http.Redirect(w, r, "/user/login", 303)
		return
	}

	http.Redirect(w, r, "/", 303)
}

func (s *Server) deleteSnippet(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("hash")
	currentUser := getAuthUserFromRequest(r)

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if currentUser.LogoutHash == hash {
		if err := s.snippetStore.Delete(int64(id), currentUser.ID); err == models.ErrNoRecord {
			http.NotFound(w, r)
			return
		} else if err != nil {
			s.serverError(w, err)
			return
		}
	} else {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/", 303)
}

func (s *Server) createSnippet(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "create", &templateData{Title: "Create snippet", CSRFField: csrf.TemplateField(r)})
}

func (s *Server) createPOST(w http.ResponseWriter, r *http.Request) {

	sForm := &snippetForm{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Expire:  r.FormValue("expire"),
		Type:    r.FormValue("type"),
	}

	errors := validation.ValidateStruct(sForm,
		validation.Field(&sForm.Title, validation.Required),
		validation.Field(&sForm.Content, validation.Required),
		validation.Field(&sForm.Expire, validation.Required, validation.By(validateInteger)),
		validation.Field(&sForm.Type, validation.Required, validation.In("Public", "Private")),
	)

	if errors != nil {
		errMap := errors.(validation.Errors)
		s.render(
			w, r,
			"create",
			&templateData{
				Title:       "Create snippet",
				Errors:      errMap,
				FormSnippet: sForm,
				CSRFField:   csrf.TemplateField(r)})

		return
	}

	currentUser := getAuthUserFromRequest(r)

	expire, err := strconv.Atoi(sForm.Expire)

	if err != nil {
		s.serverError(w, err)
		return
	}

	snippetType := true
	if sForm.Type == "Private" {
		snippetType = false
	}

	_, err = s.snippetStore.Insert(sForm.Title, sForm.Content, expire, snippetType, currentUser.ID)

	if err != nil {
		s.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/snippets", 303)
}

func (s *Server) editSnippet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	snippet, err := s.snippetStore.Get(int64(id))

	if err != nil {
		if err == models.ErrNoRecord {
			http.NotFound(w, r)
		} else {
			s.serverError(w, err)
		}
		return
	}

	currentUser := getAuthUserFromRequest(r)

	if snippet.OwnerID != currentUser.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	snippetType := "Private"
	if snippet.IsPublic {
		snippetType = "Public"
	}

	diffDate := snippet.Expires.Sub(snippet.Created).Hours() / 24

	sForm := &snippetForm{
		Title:   snippet.Title,
		Content: snippet.Content,
		Expire:  fmt.Sprintf("%d", int(diffDate)),
		Type:    snippetType,
	}

	s.render(
		w, r,
		"create",
		&templateData{
			IsEdit:      true,
			Title:       "Edit snippet",
			FormSnippet: sForm,
			CSRFField:   csrf.TemplateField(r)})
}
