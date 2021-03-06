package mysql

import (
	"database/sql"
	"testing"

	"githib.com/VladimirStepanov/snippetbox/pkg/common"
)

var (
	dsnString = common.GetEnvVariableString("DSN", "root:123@tcp(127.0.0.1:3308)/snippetbox_test?parseTime=true")
)

// GetDB ...
func GetDB(t *testing.T, dsn string) (*sql.DB, func(tables ...string)) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		for _, name := range tables {
			_, err := db.Exec("DELETE FROM  " + name)
			if err != nil {
				t.Fatal(err)
			}
		}
		defer db.Close()
	}

}
