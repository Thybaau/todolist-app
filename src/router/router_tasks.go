package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Thybaau/todolist-app/database"
	"github.com/gorilla/mux"
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
			ID:      0, //Useless because we will not use this element
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

func (s *server) handleTaskList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := s.DB.GetTaskList()
		if err != nil {
			log.Printf("Cannot load tasks, err =%v\n", err)
			http.Error(w, "Cannot load tasks", http.StatusBadRequest)
		}
		var resp = make([]jsonTask, len(tasks))
		for i, t := range tasks {
			resp[i] = jsonTask{
				ID:      t.ID,
				Content: t.Content,
				State:   t.State,
			}
		}
		// Write response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Printf("Cannot format json, err =%v\n", err)
		}
	}

}
func (s *server) handleTaskDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraire l'ID de la requête
		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de tâche invalide", http.StatusBadRequest)
		}
		//Delete Task
		err = s.DB.DeleteTask(taskID)
		if err != nil {
			log.Printf("Cannot delete task, err =%v\n", err)
			http.Error(w, "Cannot delete task", http.StatusBadRequest)
		}
		// Write response
		w.WriteHeader(http.StatusNoContent)
	}

}
