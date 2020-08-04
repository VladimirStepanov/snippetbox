package main

import (
	"fmt"

	"githib.com/VladimirStepanov/snippetbox/pkg/common"
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

	log, err := getLogger(common.GetEnvVariableString("LOG_LEVEL", "INFO"))
	if err != nil {
		return nil, err
	}

	addr := common.GetEnvVariableString("ADDR", "0.0.0.0")
	port := common.GetEnvVariableString("PORT", "8080")

	return &Config{
		addr:         fmt.Sprintf("%s:%s", addr, port),
		log:          log,
		sessionStore: sessions.NewCookieStore([]byte(common.GetEnvVariableString("SESSION_KEY", "session_key"))),
		csrfKey:      common.GetEnvVariableString("CSRF_KEY", "csrf_key"),
		dsn:          common.GetEnvVariableString("DSN", "root:123@/snippetbox?parseTime=true"),
	}, nil

}
