package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router(c Context) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", addMiddleware(indexHandler, c, isAuthed)).Methods("GET")
	r.HandleFunc("/login", addMiddleware(loginHandler, c)).Methods("GET", "POST")
	return r
}

type contextHandler func(http.ResponseWriter, *http.Request, Context)
type middleware func(contextHandler) contextHandler

func indexHandler(w http.ResponseWriter, r *http.Request, c Context) {
	w.Write([]byte("Under construction"))
}

func loginHandler(w http.ResponseWriter, r *http.Request, c Context) {
	w.Write([]byte("Login under construction"))
}

func addMiddleware(ch contextHandler, c Context, ms ...middleware) http.HandlerFunc {
	for i := len(ms) - 1; i >= 0; i-- {
		ch = ms[i](ch)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		c.Vars = mux.Vars(r)
		ch(w, r, c)
	}
}

func isAuthed(ch contextHandler) contextHandler {
	outFunc := func(w http.ResponseWriter, r *http.Request, c Context) {
		session, err := c.Store.Get(r, c.SessionName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if session.Values["user"] == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		ch(w, r, c)
	}
	return outFunc
}
