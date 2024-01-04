package router

import (
	"github.com/Thybaau/todolist-app/database"
	"github.com/gorilla/mux"
)

type server struct {
	Router *mux.Router
	DB     database.DBStore
}

func NewServer() *server {
	s := &server{
		Router: mux.NewRouter(),
	}
	s.router()
	return s
}
