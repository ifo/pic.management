package main

import (
	"html/template"
)

type Templates struct {
	Index   *template.Template
	Login   *template.Template
	NewUser *template.Template
}

func SetupTemplates(dir string) (*Templates, error) {
	index, err := template.ParseFiles(dir + "/index.html")
	if err != nil {
		return nil, err
	}
	login, err := template.ParseFiles(dir + "/login.html")
	if err != nil {
		return nil, err
	}
	newUser, err := template.ParseFiles(dir + "/newuser.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		Index:   index,
		Login:   login,
		NewUser: newUser,
	}, nil
}
