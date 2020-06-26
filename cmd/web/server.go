package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Server apllication struct
type Server struct {
}

//Routes return mux.Router with filled routes
func (s *Server) routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", s.home).Methods("GET")

	return r
}

//Start listen and serve
func (s *Server) Start() error {

	return http.ListenAndServe(":8080", s.routes())
}

//New return new Server instance
func New() *Server {
	return &Server{}
}
