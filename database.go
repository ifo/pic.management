package main

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

// TODO? make these functions DBConfig methods
func Connect(c DBConfig) (*sql.DB, error) {
	return sql.Open(c.Type, c.URL)
}

type PreparedStatements struct {
	GetUser *sql.Stmt
	NewUser *sql.Stmt
}

func SetupStmts(db *sql.DB, c DBConfig) (*PreparedStatements, error) {
	var err error
	stmts := &PreparedStatements{}
	stmts.GetUser, err = CreateGetUserQuery(db, c)
	if err != nil {
		return nil, err
	}
	stmts.NewUser, err = CreateNewUserQuery(db, c)
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func CreateGetUserQuery(db *sql.DB, c DBConfig) (*sql.Stmt, error) {
	stmt := fmt.Sprintf("SELECT id, email, password FROM %s WHERE email = $1;", c.UserTableName)
	return db.Prepare(stmt)
}

// CreateNewUserQuery returns the ID of the newly created user
func CreateNewUserQuery(db *sql.DB, c DBConfig) (*sql.Stmt, error) {
	stmt := fmt.Sprintf("INSERT INTO %s (email, password) VALUES ($1, $2); SELECT id FROM %s WHERE email = $1;",
		c.UserTableName)
	return db.Prepare(stmt)
}
