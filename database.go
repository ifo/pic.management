package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB(location, tableName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", location)
	if err != nil {
		return nil, err
	}

	// create the table
	createDBStmt := fmt.Sprintf(
		"CREATE TABLE %s (slug TEXT NOT NULL PRIMARY KEY, url TEXT);", tableName)
	_, err = db.Exec(createDBStmt)
	// but ignore existence error
	existenceError := fmt.Sprintf("table %s already exists", tableName)
	if err != nil && err.Error() != existenceError {
		return nil, err
	}

	return db, nil
}
