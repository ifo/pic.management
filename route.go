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
	r.HandleFunc("/login", addMiddleware(loginPageHandler, c)).Methods("GET")
	r.HandleFunc("/login", addMiddleware(loginHandler, c)).Methods("POST")
	r.HandleFunc("/logout", addMiddleware(logoutHandler, c)).Methods("GET")
	r.HandleFunc("/signup", addMiddleware(signupPageHandler, c)).Methods("GET")
	r.HandleFunc("/signup", addMiddleware(signupHandler, c)).Methods("POST")
	return r
}

type contextHandler func(http.ResponseWriter, *http.Request, Context)
type middleware func(contextHandler) contextHandler

func indexHandler(w http.ResponseWriter, r *http.Request, c Context) {
	c.Templates.Index.Execute(w, "")
}

func loginPageHandler(w http.ResponseWriter, r *http.Request, c Context) {
	c.Templates.Login.Execute(w, "")
}

func loginHandler(w http.ResponseWriter, r *http.Request, c Context) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		password = ""
		http.Error(w, "please provide an email and password", http.StatusBadRequest)
		return
	}

	if len(password) < 7 {
		password = ""
		http.Error(w, "passwords must be longer than 7 characters", http.StatusBadRequest)
		return
	}

	var userID int64
	var bcryptPass string
	err := c.PS.GetUser.QueryRow(email).Scan(&userID, &email, &bcryptPass)
	if err != nil && err != sql.ErrNoRows {
		password = ""
		http.Error(w, err.Error(), 500)
		return
	}

	if err == sql.ErrNoRows {
		password, bcryptPass = "", ""
		http.Error(w, "that email and password combination do not match", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(bcryptPass), []byte(password))
	password, bcryptPass = "", ""
	if err != nil {
		http.Error(w, "that email and password combination do not match", http.StatusBadRequest)
		return
	}

	// login
	session, err := c.Store.Get(r, c.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user := &User{ID: userID, Email: email}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", 303)
}

func logoutHandler(w http.ResponseWriter, r *http.Request, c Context) {
	session, err := c.Store.Get(r, c.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	session.Values["user"] = &User{}
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func signupPageHandler(w http.ResponseWriter, r *http.Request, c Context) {
	c.Templates.Signup.Execute(w, "")
}

func signupHandler(w http.ResponseWriter, r *http.Request, c Context) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordCheck := r.FormValue("password2")

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

	// TODO? add this constraint to configuration
	if len(password) < 7 {
		password, passwordCheck = "", ""
		http.Error(w, "password must be 8 characters or longer", http.StatusBadRequest)
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

	signupQuery := c.PS.Signup.QueryRow(email, string(bcryptPass))
	zero(bcryptPass)
	var userID int64
	err = signupQuery.Scan(&userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// login
	session, err := c.Store.Get(r, c.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user := &User{ID: userID, Email: email}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
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
