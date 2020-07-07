package main

import (
	"fmt"
	"net/http"
	"strconv"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
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
	s.render(w, "snippets", &templateData{Title: "Home", Snippets: snippets})
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

	s.render(w, "snippet", &templateData{Snippet: snippet})
}
