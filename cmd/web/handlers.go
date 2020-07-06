package main

import (
	"fmt"
	"net/http"
	"strconv"
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
