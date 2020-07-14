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

	if !snippet.IsPublic {
		currentUser := getAuthUserFromRequest(r)
		if currentUser == nil || currentUser.ID != snippet.OwnerID {
			http.NotFound(w, r)
			return
		}
	}

	s.render(w, r, "snippet", &templateData{Snippet: snippet})
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

	if currentUser != nil {
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
	}

	http.Redirect(w, r, "/", 303)
}
