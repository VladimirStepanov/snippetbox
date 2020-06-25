package main

import "github.com/gorilla/mux"

//Routes return mux.Router with filled routes
func (s *Server) Routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", s.home).Methods("GET")

	return r
}
