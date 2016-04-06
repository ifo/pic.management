package setup

import (
	"database/sql"
	"fmt"

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

func CreateDBTables(db *sql.DB, c DBConfig) error {
	// TODO setup all the tables

	// create user table
	createUserTable := fmt.Sprintf(
		"CREATE TABLE %s (id INTEGER PRIMARY KEY ASC, email TEXT UNIQUE, password TEXT);",
		c.UserTableName)
	_, err := db.Exec(createUserTable)
	// ignore existence error
	if err != nil && err.Error() != existenceError(c.UserTableName) {
		return err
	}

	return nil
}

func existenceError(name string) string {
	return "table " + name + " already exists"
}
