package main

import (
	"fmt"

	"githib.com/VladimirStepanov/snippetbox/pkg/models/mysql"
)

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
		config.addr,
		config.log,
		&mysql.UsersStore{DB: db},
		&mysql.SnippetStore{DB: db},
		config.sessionStore,
		config.csrfKey,
	)

	if err = serv.Start(); err != nil {
		config.log.Errorf("Error while Start server... %v\n", err)
	}

}
