package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server apllication struct
type Server struct {
	addr          string
	log           *logrus.Logger
	templateCache map[string]*template.Template
}

//Routes return mux.Router with filled routes
func (s *Server) routes() http.Handler {
	r := mux.NewRouter()

	strPref := http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/")))
	r.PathPrefix("/static/").Handler(strPref)
	r.HandleFunc("/", s.home).Methods("GET")
	return s.loggerMiddleware(r)
}

//Start listen and serve
func (s *Server) Start() error {

	templateCache, err := newTemplateCache("./ui/html")

	if err != nil {
		return err
	}

	s.templateCache = templateCache

	srv := &http.Server{
		Handler:      s.routes(),
		Addr:         s.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.log.Infof("Server start at addr %s\n", s.addr)

	return srv.ListenAndServe()
}

//New return new Server instance
func New(addr string, log *logrus.Logger) *Server {
	return &Server{
		addr: addr,
		log:  log,
	}
}
