package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// TODO use sessions
var store *sessions.CookieStore

func main() {
	var (
		sessionSecret = flag.String("session", os.Getenv("SESSION_SECRET"), "Set the session secret")
	)
	flag.Parse()
	store = sessions.NewCookieStore([]byte(*sessionSecret))

	c := Context{}
	r := router(c)
	http.ListenAndServe(":3000", r)
}

type Context struct {
	Vars map[string]string
	// TODO have db connection or other data here
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
