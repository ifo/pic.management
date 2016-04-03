package main

import (
	"log"
	"net/http"
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
