package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(128),
comment TEXT,
repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date)
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool

	if os.IsNotExist(err) {
		// Файл точно не существует -> надо ставить схему
		install = true

	} else if err != nil {
		return err
	}

	var openErr error
	db, openErr = sql.Open("sqlite", dbFile)

	if openErr != nil {
		return openErr
	}

	if install == true {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
