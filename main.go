package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
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
		defaultDbURL         = "file:data/sqlite.db"
		defaultPort          = 3000
		sessionSecret        = flag.String("session", defaultSessionSecret, "Set the session secret")
		dbURL                = flag.String("db", defaultDbURL, "Set the database connection string")
		port                 = flag.Int("port", defaultPort, "Set the server port")
		envSessionSecret     = os.Getenv("SESSION_SECRET")
		envDbURL             = os.Getenv("DATABASE_URL")
		outPort              = os.Getenv("PORT")
	)
	flag.Parse()

	if *sessionSecret == defaultSessionSecret && envSessionSecret != "" {
		sessionSecret = &envSessionSecret
	}
	if *dbURL == defaultDbURL && envDbURL != "" {
		dbURL = &envDbURL
	}
	if *port != defaultPort || outPort == "" {
		outPort = strconv.Itoa(defaultPort)
	}

	store := sessions.NewCookieStore([]byte(*sessionSecret))
	db, err := SetupDB(*dbURL)
	if err != nil {
		return nil, "", err
	}
	return &Context{DB: db, Store: store}, outPort, nil
}

type Context struct {
	Vars  map[string]string
	DB    *sql.DB
	Store *sessions.CookieStore
}

func router(c Context) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", injectContext(indexHandler, c)).Methods("GET")
	r.HandleFunc("/login", injectContext(loginHandler, c)).Methods("POST")
	return r
}

func injectContext(fn func(http.ResponseWriter, *http.Request, Context), c Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Vars = mux.Vars(r)
		fn(w, r, c)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, c Context) {
	w.Write([]byte("Under construction"))
}

func loginHandler(w http.ResponseWriter, r *http.Request, c Context) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
