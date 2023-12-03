package main

import (
	"log"
	"net/http"

	"github.com/Thybaau/todolist-app/router"
)

func main() {
	log.Printf("Todo-list app Golang")
	r := router.Router()

	// Server connexion
	// http.HandleFunc("/", srv.serveHTTP)
	log.Printf("Serving HTTP on port 9000")
	err := http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatal(err)
	}
}
