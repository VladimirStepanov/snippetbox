package main

import (
	"net/http"
)

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
