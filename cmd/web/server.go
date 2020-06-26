package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Server apllication struct
type Server struct {
	addr string
}

//Routes return mux.Router with filled routes
func (s *Server) routes() *mux.Router {
	r := mux.NewRouter()

	strPref := http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/")))
	r.PathPrefix("/static/").Handler(strPref)
	r.HandleFunc("/", s.home).Methods("GET")
	return r
}

//Start listen and serve
func (s *Server) Start() error {

	srv := &http.Server{
		Handler:      s.routes(),
		Addr:         s.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

//New return new Server instance
func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}
