package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Thybaau/todolist-app/database"
)

type jsonTask struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	State   bool   `json:"state"`
}

func (s *server) handleTaskCreate() http.HandlerFunc {
	type request struct {
		Content string `json:"content"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Cannot parse task body. err=%v\n", err)
			http.Error(w, "Cannot parse task body from json", http.StatusBadRequest)
		}

		// Insert task in database
		t := &database.Task{
			Content: req.Content,
			State:   false,
		}
		id, err := s.DB.CreateTask(t)
		if err != nil {
			log.Printf("Cannot create task in database. err=%v\n", err)
			http.Error(w, "Cannot create task in database", http.StatusBadRequest)
		}

		// Write response
		var resp = jsonTask{
			ID:      id,
			Content: t.Content,
			State:   t.State,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Printf("Cannot format json, err =%v\n", err)
		}
	}
}
