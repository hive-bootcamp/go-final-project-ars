package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var Db *sql.DB
var schema = `
	CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
	);

	CREATE INDEX scheduler_idx on scheduler (date);
`

func Init(dbFile string) error {
	envDBFile := os.Getenv("TODO_DBFILE")
	var finalDBFile string
	if envDBFile != "" {
		finalDBFile = envDBFile
	} else {
		finalDBFile = dbFile
	}
	_, err := os.Stat(finalDBFile)

	var install bool
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		install = true
	}

	Db, err = sql.Open("sqlite", finalDBFile)
	if err != nil {
		return err
	}
	if install {
		_, err := Db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
