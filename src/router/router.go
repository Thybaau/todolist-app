package router

func (s *server) router() {
	s.Router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.Router.HandleFunc("/tasks/list/", s.handleTaskList()).Methods("GET")
	s.Router.HandleFunc("/tasks/create/", s.handleTaskCreate()).Methods("POST")
	s.Router.HandleFunc("/tasks/{id:[0-9]+}", s.handleTaskDelete()).Methods("DELETE")
	s.Router.HandleFunc("/tasks/{id:[0-9]+}", s.handleTaskEdit()).Methods("PUT")
}
