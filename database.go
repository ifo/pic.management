package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	UserTableName  string
	PhotoTableName string
	ImageTableName string
	// TODO? add Album and Group
}

func SetupDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	// TODO setup the tables

	return db, nil
}
