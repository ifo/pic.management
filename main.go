package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

func main() {
	context, port, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	r := router(*context)
	log.Println("Starting server on port " + port)
	http.ListenAndServe(":"+port, r)
}

func setup() (*Context, string, error) {
	var (
		defaultSessionSecret = "session-secret"
		defaultDBURL         = "file:db/sqlite.db"
		defaultPort          = 3000
		defaultSessionName   = "session-name"
		sessionSecret        = flag.String("session", defaultSessionSecret, "Set the session secret")
		dbURL                = flag.String("db", defaultDBURL, "Set the database connection string")
		port                 = flag.Int("port", defaultPort, "Set the server port")
		sessionName          = flag.String("session-name", defaultSessionName, "Set the session name")
		envSessionSecret     = os.Getenv("SESSION_SECRET")
		envDBURL             = os.Getenv("DATABASE_URL")
		outPort              = os.Getenv("PORT")
		envSessionName       = os.Getenv("SESSION_NAME")
	)
	flag.Parse()

	if *sessionSecret == defaultSessionSecret && envSessionSecret != "" {
		sessionSecret = &envSessionSecret
	}
	if *dbURL == defaultDBURL && envDBURL != "" {
		dbURL = &envDBURL
	}
	if *port != defaultPort || outPort == "" {
		outPort = strconv.Itoa(defaultPort)
	}
	if *sessionName == defaultSessionName && envSessionName != "" {
		sessionName = &envSessionName
	}

	store := sessions.NewCookieStore([]byte(*sessionSecret))

	dbType, err := getDBType(*dbURL)
	if err != nil {
		return nil, "", err
	}

	dbConfig := DBConfig{
		URL:  *dbURL,
		Type: dbType,
	}
	db, err := SetupDB(dbConfig)
	if err != nil {
		return nil, "", err
	}
	return &Context{DB: db, Store: store, SessionName: *sessionName}, outPort, nil
}

type Context struct {
	Vars        map[string]string
	DB          *sql.DB
	Store       *sessions.CookieStore
	SessionName string
}

func getDBType(url string) (string, error) {
	proto := getProtocol(url)
	switch proto {
	case "file", "sqlite", "sqlite3":
		return "sqlite3", nil
	case "postgres":
		return "postgres", nil
	default:
		return "", errors.New("Unknown database protocol")
	}
}

func getProtocol(url string) string {
	for i, r := range url {
		if r == ':' {
			return url[:i]
		}
	}
	return ""
}
