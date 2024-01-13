package router

func (s *server) router() {
	s.Router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.Router.HandleFunc("/tasks", s.handleTaskList()).Methods("GET")
	s.Router.HandleFunc("/tasks", s.handleTaskCreate()).Methods("POST")
	s.Router.HandleFunc("/tasks/{id:[0-9]+}", s.handleTaskDelete()).Methods("DELETE")
	s.Router.HandleFunc("/tasks/{id:[0-9]+}", s.handleTaskEdit()).Methods("PUT")
	s.Router.HandleFunc("/tasks/state/{id:[0-9]+}", s.handleTaskState()).Methods("PUT")
}
