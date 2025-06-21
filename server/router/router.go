package router

import "github.com/Thybaau/todolist-app/middleware"

func (s *server) router() {
	taskRouter := s.Router.PathPrefix("/tasks").Subrouter()
	taskRouter.Use(middleware.JWTMiddleware)

	taskRouter.HandleFunc("", s.handleTaskList()).Methods("GET")
	taskRouter.HandleFunc("", s.handleTaskCreate()).Methods("POST")
	taskRouter.HandleFunc("/{id:[0-9]+}", s.handleTaskDelete()).Methods("DELETE")
	taskRouter.HandleFunc("/{id:[0-9]+}", s.handleTaskEdit()).Methods("PUT")
	taskRouter.HandleFunc("/state/{id:[0-9]+}", s.handleTaskState()).Methods("PUT")

	s.Router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.Router.HandleFunc("/account/signup", s.handleSignUp()).Methods("POST")
	s.Router.HandleFunc("/account/login", s.handleLogin()).Methods("POST")
}
