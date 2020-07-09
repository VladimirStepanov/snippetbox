package main

import (
	"database/sql"
	"net/http"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
func (s *Server) serverError(w http.ResponseWriter, err error) {
	s.log.Errorf("Internal error: %v %s", err, string(debug.Stack()))
	http.Error(w, "Internal error", http.StatusInternalServerError)
}

func (s *Server) addFlashMessage(w http.ResponseWriter, r *http.Request, message string) error {
	session, err := s.session.Get(r, "flash")
	if err != nil {
		return err
	}

	session.AddFlash(message)

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil

}

func (s *Server) getFlashes(w http.ResponseWriter, r *http.Request) ([]interface{}, error) {
	session, err := s.session.Get(r, "flash")

	if err != nil {
		return nil, err
	}

	flashes := session.Flashes()

	err = session.Save(r, w)

	if err != nil {
		s.serverError(w, err)
		return nil, err
	}

	return flashes, nil

}
