package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Thybaau/todolist-app/database"
	"github.com/Thybaau/todolist-app/middleware"
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
		//Decode and check fields in request
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot decode task body from json", http.StatusBadRequest, err)
			return
		}
		if req.Content == "" {
			middleware.NewHTTPError(w, "Key 'content' cannot be empty", http.StatusForbidden, nil)
			return
		}

		// Insert task in database
		t := &database.Task{
			ID:      0, //Useless because we will not use this element
			Content: req.Content,
			State:   false,
		}
		id, err := s.DB.CreateTask(t)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot create task in database", http.StatusBadRequest, err)
			return
		}

		// Write response
		var resp = jsonTask{
			ID:      id,
			Content: t.Content,
			State:   t.State,
		}
		middleware.JSONResponse(w, http.StatusOK, resp)
	}
}

func (s *server) handleTaskList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp interface{}
		queryParams := r.URL.Query()

		// If we did not put any query parameter, we get all the task list
		if len(queryParams) == 0 {
			tasks, err := s.DB.GetTaskList()
			if err != nil {
				middleware.NewHTTPError(w, "Cannot load tasks", http.StatusInternalServerError, err)
			}
			resp = make([]jsonTask, len(tasks))
			for i, t := range tasks {
				resp.([]jsonTask)[i] = jsonTask{
					ID:      t.ID,
					Content: t.Content,
					State:   t.State,
				}
			}
			// If we put query parameter 'id', we get task with this id
		} else {
			taskID := queryParams.Get("id")
			if taskID == "" {
				middleware.NewHTTPError(w, "Query parameter 'id' not found", http.StatusNotFound, nil)
				return
			}
			ID, _ := strconv.Atoi(taskID)
			task, err := s.DB.GetTask(ID)
			if err != nil {
				message := fmt.Sprintf("Task id=%v not found", taskID)
				middleware.NewHTTPError(w, message, http.StatusNotFound, err)
				return
			}
			resp = jsonTask{
				ID:      task.ID,
				Content: task.Content,
				State:   task.State,
			}
		}
		// Write response
		middleware.JSONResponse(w, http.StatusOK, resp)
	}
}

func (s *server) handleTaskDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract request ID
		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["id"])
		if err != nil {
			middleware.NewHTTPError(w, "Invalid task ID", http.StatusBadRequest, err)
			return
		}

		//Delete Task
		err = s.DB.DeleteTask(taskID)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot delete task", http.StatusBadRequest, err)
			return
		}

		// Write response
		successMessage := fmt.Sprintf("successfully deleted task with id=%v", taskID)
		jsonResp := map[string]string{"message": successMessage}
		middleware.JSONResponse(w, http.StatusOK, jsonResp)
	}
}

func (s *server) handleTaskEdit() http.HandlerFunc {
	type request struct {
		Content string `json:"content"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode RequestBody
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot parse task body", http.StatusBadRequest, err)
			return
		}
		//Valide fields in the request
		if req.Content == "" {
			middleware.NewHTTPError(w, "Key 'content' cannot be empty", http.StatusForbidden, nil)
			return
		}

		// Extract ID from path parameter
		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["id"])
		if err != nil {
			middleware.NewHTTPError(w, "Invalid ID", http.StatusBadRequest, err)
			return
		}
		err = s.DB.EditTask(taskID, req.Content)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot edit task", http.StatusBadRequest, err)
			return
		}

		// Write response
		task, err := s.DB.GetTask(taskID)
		if err != nil {
			message := fmt.Sprintf("Cannot get task informations for id=%v", taskID)
			middleware.NewHTTPError(w, message, http.StatusBadRequest, err)
			return
		}
		var resp = jsonTask{
			ID:      task.ID,
			Content: task.Content,
			State:   task.State,
		}
		middleware.JSONResponse(w, http.StatusOK, resp)
	}
}

func (s *server) handleTaskState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract request ID
		vars := mux.Vars(r)
		taskID, err := strconv.Atoi(vars["id"])
		if err != nil {
			middleware.NewHTTPError(w, "Invalid task ID", http.StatusBadRequest, err)
			return
		}

		task, err := s.DB.ChangeTaskState(taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				middleware.NewHTTPError(w, "Task not found", http.StatusNotFound, err)
				return
			} else {
				middleware.NewHTTPError(w, "Cannot change task state", http.StatusBadRequest, err)
				return
			}
		}

		// Write response
		var resp = jsonTask{
			ID:      task.ID,
			Content: task.Content,
			State:   task.State,
		}
		middleware.JSONResponse(w, http.StatusOK, resp)

	}
}
