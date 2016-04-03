package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	URL            string
	Type           string
	UserTableName  string
	PhotoTableName string
	ImageTableName string
	// TODO? add Album and Group
}

func SetupDB(c DBConfig) (*sql.DB, error) {
	db, err := sql.Open(c.Type, c.URL)
	if err != nil {
		return nil, err
	}

	// TODO setup the tables

	return db, nil
}
