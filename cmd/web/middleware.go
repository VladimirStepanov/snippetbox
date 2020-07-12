package main

import (
	"context"
	"net/http"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
)

type contextKey string

var contextKeyUser = contextKey("user")

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *Server) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			lwr := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lwr, r)
			s.log.Infof("%s %s %s %s %d", r.Method, r.Proto, r.RemoteAddr, r.RequestURI, lwr.statusCode)
		})

}

func (s *Server) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			session, err := s.session.Get(r, "SID")
			if err != nil {
				removeSession(w, r, session)
				next.ServeHTTP(w, r)
				return
			}
			if len(session.Values) == 2 {
				userID := session.Values["userID"].(int64)

				u, err := s.userStore.Get(userID)
				if err == models.ErrNoRecord {
					removeSession(w, r, session)
					next.ServeHTTP(w, r)
					return
				} else if err != nil {
					s.serverError(w, err)
					return
				}
				u.LogoutHash = session.Values["logoutHash"].(string)
				ctx := context.WithValue(r.Context(), contextKeyUser, u)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
}

func (s *Server) accessOnlyNotAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if getAuthUserFromRequest(r) != nil {
				http.Redirect(w, r, "/", 303)
				return
			}
			next.ServeHTTP(w, r)
		})
}
