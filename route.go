package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func router(c Context) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", addMiddleware(indexHandler, c, isAuthed)).Methods("GET")
	r.HandleFunc("/login", addMiddleware(loginHandler, c)).Methods("GET", "POST")
	r.HandleFunc("/newuser", addMiddleware(newUserHandler, c)).Methods("POST")
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

func newUserHandler(w http.ResponseWriter, r *http.Request, c Context) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordCheck := r.FormValue("repeat")

	// TODO for all errors redirect and show flash message
	if email == "" {
		password, passwordCheck = "", ""
		http.Error(w, "email not provided", http.StatusBadRequest)
		return
	}

	if password != passwordCheck {
		password, passwordCheck = "", ""
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	// TODO handle database connection error
	if _, err := c.PS.GetUser.Query(email); err != nil && err != sql.ErrNoRows {
		password, passwordCheck = "", ""
		http.Error(w, "email is already taken", http.StatusBadRequest)
		return
	}

	// save user
	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password, passwordCheck = "", ""
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUserQuery := c.PS.NewUser.QueryRow(email, bcryptPass)
	zero(bcryptPass)
	var userID int64
	err = newUserQuery.Scan(&userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// login
	session, err := c.Store.Get(r, c.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	user := &User{ID: userID, Email: email}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	http.Redirect(w, r, "/", 303)
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

func zero(bts []byte) {
	for i := range bts {
		bts[i] = 0
	}
}
