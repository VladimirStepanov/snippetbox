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

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	var err error
	page := 1
	pageStr := r.URL.Query().Get("page")

	if pageStr != "" {
		retPage, err := strconv.Atoi(pageStr)
		if err != nil {
			//Add flash message
			s.serverError(w, err)
			return
		}

		if retPage < 1 {
			s.serverError(w, fmt.Errorf("Page must be greater than zero"))
			return
		}

		page = retPage
	}
	snippets, err := s.snippetStore.LatestAll(-1, 10, page)

	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, r, "snippets", &templateData{Title: "Home", Snippets: snippets})
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
		// check for user session
		http.NotFound(w, r)
		return
	}

	s.render(w, r, "snippet", &templateData{Snippet: snippet})
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "signup", &templateData{CSRFField: csrf.TemplateField(r)})
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
	http.Redirect(w, r, "/", 303)

}
