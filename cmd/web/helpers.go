package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
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

func (s *Server) addNewUserSession(w http.ResponseWriter, r *http.Request, id int64) error {
	session, err := s.session.Get(r, "SID")

	if err != nil {
		return err
	}

	hasher := md5.New()

	_, err = hasher.Write([]byte(fmt.Sprintf("%d%s", id, time.Now().String())))

	if err != nil {
		return err
	}

	session.Values["userID"] = id
	session.Values["logoutHash"] = hex.EncodeToString(hasher.Sum(nil))

	if err = session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func removeSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func getAuthUserFromRequest(r *http.Request) *models.User {
	u, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}

	return u
}
