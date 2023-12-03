package router

import "github.com/gorilla/mux"

func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex()).Methods("GET")
	return router
}
