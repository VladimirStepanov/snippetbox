package main

import (
	"fmt"
	"net/http"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "My path is "+r.URL.Path)
}
