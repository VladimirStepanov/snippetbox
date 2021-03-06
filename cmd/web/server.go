package main

import (
	"html/template"
	"net/http"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

//Server apllication struct
type Server struct {
	addr          string
	log           *logrus.Logger
	templateCache map[string]*template.Template
	userStore     models.UserRepository
	snippetStore  models.SnippetRepository
	session       *sessions.CookieStore
	csrfKey       string
}

//Routes return mux.Router with filled routes
func (s *Server) routes() http.Handler {

	CSRF := csrf.Protect([]byte(s.csrfKey), csrf.Secure(false))

	r := mux.NewRouter()

	strPref := http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/")))
	r.PathPrefix("/static/").Handler(strPref)
	r.HandleFunc("/", s.home).Methods("GET")
	r.Handle("/snippets", s.accessOnlyAuth(http.HandlerFunc(s.userSnippets))).Methods("GET")
	r.Handle("/snippet/create", s.accessOnlyAuth(http.HandlerFunc(s.createSnippet))).Methods("GET")
	r.Handle("/snippet/create", s.accessOnlyAuth(http.HandlerFunc(s.createPOST))).Methods("POST")
	r.Handle("/snippet/delete/{id:[0-9]+}", s.accessOnlyAuth(http.HandlerFunc(s.deleteSnippet))).Methods("GET")
	r.Handle("/snippet/edit/{id:[0-9]+}", s.accessOnlyAuth(http.HandlerFunc(s.editSnippet))).Methods("GET")
	r.Handle("/snippet/edit/{id:[0-9]+}", s.accessOnlyAuth(http.HandlerFunc(s.editPOST))).Methods("POST")
	r.HandleFunc("/snippet/{id:[0-9]+}", s.showSnippet).Methods("GET")
	r.Handle("/user/signup", s.accessOnlyNotAuth(http.HandlerFunc(s.signUpPOST))).Methods("POST")
	r.Handle("/user/signup", s.accessOnlyNotAuth(http.HandlerFunc(s.signUp))).Methods("GET")
	r.Handle("/user/login", s.accessOnlyNotAuth(http.HandlerFunc(s.showLogin))).Methods("GET")
	r.Handle("/user/login", s.accessOnlyNotAuth(http.HandlerFunc(s.loginPOST))).Methods("POST")
	r.Handle("/user/logout", s.accessOnlyAuth(http.HandlerFunc(s.logout))).Methods("GET")
	return s.loggerMiddleware(s.authUser(CSRF(r)))
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
func New(
	config *Config,
	ur models.UserRepository,
	sr models.SnippetRepository,
) *Server {

	return &Server{
		addr:         config.addr,
		log:          config.log,
		userStore:    ur,
		snippetStore: sr,
		session:      config.sessionStore,
		csrfKey:      config.csrfKey,
	}
}
