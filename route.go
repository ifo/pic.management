package main

import (
	"net/http"

	"github.com/gorilla/mux"
	// TODO use sessions
	//"github.com/gorilla/sessions"
)

func router(c Context) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", injectContext(indexHandler, c)).Methods("GET")
	r.HandleFunc("/login", injectContext(loginHandler, c)).Methods("POST")
	return r
}

type contextHandler func(http.ResponseWriter, *http.Request, Context)

func injectContext(fn contextHandler, c Context) http.HandlerFunc {
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
