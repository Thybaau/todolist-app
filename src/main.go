package main

import (
	"log"
	"net/http"

	"github.com/Thybaau/todolist-app/database"
	"github.com/Thybaau/todolist-app/middleware"
	"github.com/Thybaau/todolist-app/router"
)

const (
	host     = "database"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "todolist_db"
)

func main() {
	log.Printf("Running todo-list app Golang...")
	srv := router.NewServer()

	// Database connexion
	srv.DB = database.DBStore{}
	err := srv.DB.Connect(host, port, user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to database")
	defer srv.DB.Close()

	// Server connexion
	// http.HandleFunc("/", srv.serveHTTP)
	srv.Router.Use(middleware.LogRequests)
	log.Printf("Running server on port 9000")
	err = http.ListenAndServe(":9000", srv.Router)
	if err != nil {
		log.Fatal(err)
	}
}
