package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Thybaau/todolist-app/database"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Task list

func TestHandleTaskList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(1, "Task 1", false).
		AddRow(2, "Task 2", false)

	mock.ExpectQuery("SELECT id, content, state FROM tasks").WillReturnRows(rows)
	srv := &server{
		DB: &database.DBStore{DB: db},
	}
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	srv.handleTaskList()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}

	expectedResp := `[
		{
		  "id": 1,
		  "content": "Task 1",
		  "state": false
		},
		{
		  "id": 2,
		  "content": "Task 2",
		  "state": false
		}
	  ]`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleTaskListOneTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(2, "Task 2", false)

	query := "SELECT id, content, state FROM tasks WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(2).WillReturnRows(rows)
	srv := &server{
		DB: &database.DBStore{DB: db},
	}
	req := httptest.NewRequest("GET", "/tasks?id=2", nil)
	w := httptest.NewRecorder()
	srv.handleTaskList()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}

	expectedResp := `{"id":2,"content":"Task 2","state":false}`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

// Delete Task

func TestHandleTaskDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	taskID := "12"
	delete := "DELETE FROM tasks WHERE id = \\$1"
	mock.ExpectExec(delete).WithArgs(12).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("DELETE", "/tasks/"+taskID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": taskID})
	w := httptest.NewRecorder()
	srv.handleTaskDelete()(w, req)

	expectedResp := `{"message": "successfully deleted task with id=12"}`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

// Create task

func TestHandleTaskCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	task := &database.Task{
		ID:      0,
		Content: "test task content",
		State:   false,
	}

	insert := "INSERT INTO tasks (content,state) VALUES ($1, $2) RETURNING id"
	mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(task.Content, task.State).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	requestBody := []byte(`{"content": "test task content"}`)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	srv.handleTaskCreate()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
	expectedResp := `{
		"id": 1,
		"content": "test task content",
		"state": false
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

}

// Edit task

func TestHandleTaskEdit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	taskID := "12"
	content := "test task content"

	query := "SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(12).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectExec("UPDATE tasks SET content = \\$1 WHERE id = \\$2").
		WithArgs(content, 12).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(12, content, false)

	query = "SELECT id, content, state FROM tasks WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(12).WillReturnRows(rows)

	requestBody := []byte(`{"content": "test task content"}`)
	req := httptest.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": taskID})
	w := httptest.NewRecorder()
	srv.handleTaskEdit()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
	expectedResp := `{
		"id": 12,
		"content": "test task content",
		"state": false
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

// Change task state

func TestHandleChangeTaskState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()
	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	taskID := 12
	state := true

	query := "SELECT id, content, state FROM tasks WHERE id = \\$1"
	rows1 := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(taskID, "Task 1", !state)
	mock.ExpectQuery(query).WithArgs(taskID).WillReturnRows(rows1)

	mock.ExpectExec("UPDATE tasks SET state = \\$1 WHERE id = \\$2").
		WithArgs(state, taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows2 := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(taskID, "Task 1", state)
	mock.ExpectQuery(query).WithArgs(taskID).WillReturnRows(rows2)

	req := httptest.NewRequest("PUT", "/tasks/state/12", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "12"})
	w := httptest.NewRecorder()
	srv.handleTaskState()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
	expectedResp := `{
		"id": 12,
		"content": "Task 1",
		"state": true
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
