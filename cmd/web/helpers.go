package main

import (
	"net/http"
	"runtime/debug"
)

func (s *Server) serverError(w http.ResponseWriter, err error) {
	s.log.Errorf("Internal error: %v %s", err, string(debug.Stack()))
	http.Error(w, "Internal error", http.StatusInternalServerError)
}
