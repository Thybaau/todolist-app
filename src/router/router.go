package router

import "github.com/gorilla/mux"

func (s *server) router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex()).Methods("GET")
	return router
}
