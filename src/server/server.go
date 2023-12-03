package server

import (
	"github.com/Thybaau/todolist-app/database"
	"github.com/Thybaau/todolist-app/router"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	DB     database.DBStore
}

func NewServer() *Server {
	r := router.Router()
	s := &Server{
		Router: r,
	}
	return s
}
