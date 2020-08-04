package main

import (
	"fmt"

	"githib.com/VladimirStepanov/snippetbox/pkg/models/mysql"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("conf.env")
}

func main() {

	config, err := NewConfig()

	if err != nil {
		fmt.Printf("Error while create Config: %v\n", err)
	}

	db, err := openDB(config.dsn)

	if err != nil {
		config.log.Errorf("Error while open DB connection %v", err)
		return
	}

	serv := New(
		config,
		&mysql.UsersStore{DB: db},
		&mysql.SnippetStore{DB: db},
	)

	if err = serv.Start(); err != nil {
		config.log.Errorf("Error while Start server... %v\n", err)
	}

}
