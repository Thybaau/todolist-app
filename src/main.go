package main

import (
	"log"
	"net/http"

	"github.com/Thybaau/todolist-app/database"
	"github.com/Thybaau/todolist-app/router"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "todolist_db"
)

func main() {
	log.Printf("Running todo-list app Golang...")
	r := router.Router()

	// Database connexion
	db := &database.DBStore{}
	err := db.Connect(host, port, user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Server connexion
	// http.HandleFunc("/", srv.serveHTTP)
	log.Printf("Serving HTTP on port 9000")
	err = http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatal(err)
	}
}
