package router

import (
	"fmt"
	"net/http"
)

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Todo-List by Thibault")
	}
}
