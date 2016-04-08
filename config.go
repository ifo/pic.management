package main

import (
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

type Context struct {
	Vars        map[string]string
	Templates   *Templates
	PS          *PreparedStatements
	Store       *sessions.CookieStore
	SessionName string
}

func setup() (*Context, string, error) {
	var (
		defaultSessionSecret = "session-secret"
		sessionSecret        = flag.String("session", defaultSessionSecret, "Set the session secret")
		envSessionSecret     = os.Getenv("SESSION_SECRET")

		defaultDBURL = "file:db/sqlite.db"
		dbURL        = flag.String("db", defaultDBURL, "Set the database connection string")
		envDBURL     = os.Getenv("DATABASE_URL")

		defaultPort = 3000
		port        = flag.Int("port", defaultPort, "Set the server port")
		outPort     = os.Getenv("PORT")

		defaultSessionName = "session-name"
		sessionName        = flag.String("session-name", defaultSessionName, "Set the session name")
		envSessionName     = os.Getenv("SESSION_NAME")

		defaultTemplatesDir = "templates"
		templatesDir        = flag.String("templates", defaultTemplatesDir, "Set the templates directory")
		envTemplatesDir     = os.Getenv("TEMPLATES_DIR")

		defaultUserTableName = "user"
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
	if *templatesDir == defaultTemplatesDir && envTemplatesDir != "" {
		templatesDir = &envTemplatesDir
	}

	store := sessions.NewCookieStore([]byte(*sessionSecret))

	dbType, err := getDBType(*dbURL)
	if err != nil {
		return nil, "", err
	}

	dbConfig := DBConfig{
		URL:           *dbURL,
		Type:          dbType,
		UserTableName: defaultUserTableName,
	}
	db, err := Connect(dbConfig)
	if err != nil {
		return nil, "", err
	}
	stmts, err := SetupStmts(db, dbConfig)
	if err != nil {
		return nil, "", err
	}

	templates, err := SetupTemplates(*templatesDir)
	if err != nil {
		return nil, "", err
	}

	return &Context{
		Templates:   templates,
		PS:          stmts,
		Store:       store,
		SessionName: *sessionName,
	}, outPort, nil
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
