package router

func (s *server) router() {
	s.Router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.Router.HandleFunc("/tasks/create/", s.handleTaskCreate()).Methods("POST")
}
