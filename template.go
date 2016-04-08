package main

import (
	"html/template"
)

type Templates struct {
	Index *template.Template
}

func SetupTemplates(dir string) (*Templates, error) {
	index, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		Index: index,
	}, nil
}
