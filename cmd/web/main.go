package main

import (
	"flag"
	"fmt"

	"githib.com/VladimirStepanov/snippetbox/pkg/models/mysql"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

func getLogger(levelString string) (*logrus.Logger, error) {
	log := logrus.New()
	level, err := logrus.ParseLevel(levelString)

	if err != nil {
		return nil, err
	}

	log.SetLevel(level)

	return log, nil
}

func main() {
	addr := flag.String("addr", ":8080", "Listen addr")
	logLevel := flag.String("level", "INFO", "Log level")
	dsn := flag.String("dsn", "root:123@/snippetbox?parseTime=true", "Dsn")
	key := flag.String("key", "test-123", "Session key")

	flag.Parse()

	sessionStore := sessions.NewCookieStore([]byte(*key))

	log, err := getLogger(*logLevel)

	if err != nil {
		fmt.Printf("Error while parse Log level: %v\n", err)
		return
	}

	db, err := openDB(*dsn)

	if err != nil {
		log.Errorf("Error while open DB connection %v", err)
		return
	}

	serv := New(*addr, log, &mysql.UsersStore{DB: db}, &mysql.SnippetStore{DB: db}, sessionStore)

	if err = serv.Start(); err != nil {
		log.Errorf("Error while Start server... %v\n", err)
	}

}
