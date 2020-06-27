package main

import (
	"net/http"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.render(w, "snippets", &templateData{Title: "Home"})
}
