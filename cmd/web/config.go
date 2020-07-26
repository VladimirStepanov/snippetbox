package main

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

//Config struct for web application
type Config struct {
	addr         string
	log          *logrus.Logger
	sessionStore *sessions.CookieStore
	csrfKey      string
	dsn          string
}

func getEnvVariableString(key, defaultValue string) string {
	var res string
	if res = os.Getenv(key); res == "" {
		res = defaultValue
	}
	return res
}

func getLogger(levelString string) (*logrus.Logger, error) {
	log := logrus.New()
	level, err := logrus.ParseLevel(levelString)

	if err != nil {
		return nil, err
	}

	log.SetLevel(level)

	return log, nil
}

//NewConfig ...
func NewConfig() (*Config, error) {

	log, err := getLogger(getEnvVariableString("LOG_LEVEL", "INFO"))
	if err != nil {
		return nil, err
	}

	addr := getEnvVariableString("ADDR", "0.0.0.0")
	port := getEnvVariableString("PORT", "8080")

	return &Config{
		addr:         fmt.Sprintf("%s:%s", addr, port),
		log:          log,
		sessionStore: sessions.NewCookieStore([]byte(getEnvVariableString("SESSION_KEY", "session_key"))),
		csrfKey:      getEnvVariableString("CSRF_KEY", "csrf_key"),
		dsn:          getEnvVariableString("DSN", "root:123@/snippetbox?parseTime=true"),
	}, nil

}
