package router

import (
	"bytes"
	"database/sql"
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

func TestHandleTaskListIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	query := "SELECT id, content, state FROM tasks WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(2).WillReturnError(sql.ErrNoRows)
	srv := &server{
		DB: &database.DBStore{DB: db},
	}
	req := httptest.NewRequest("GET", "/tasks?id=2", nil)
	w := httptest.NewRecorder()
	srv.handleTaskList()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}

	expectedResp := `{
		"error": "Task id=2 not found",
		"detail": "sql: no rows in result set"
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleTaskListQueryNotFound(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}
	req := httptest.NewRequest("GET", "/tasks?id=", nil)
	w := httptest.NewRecorder()
	srv.handleTaskList()(w, req)

	expectedResp := `{
		"error": "Query parameter 'id' not found",
		"detail": ""
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
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

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}

	expectedResp := `{"message": "successfully deleted task with id=12"}`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleTaskDeleteIDNotFound(t *testing.T) {
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
	mock.ExpectExec(delete).WithArgs(12).WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("DELETE", "/tasks/"+taskID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": taskID})
	w := httptest.NewRecorder()
	srv.handleTaskDelete()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}

	expectedResp := `{
		"error": "Cannot delete task",
		"detail": "task with ID 12 does not exist"
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
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

func TestHandleTaskCreateContentEmpty(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	requestBody := []byte(`{"content": ""}`)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	srv.handleTaskCreate()(w, req)

	expectedResp := `{
		"error": "Key 'content' cannot be empty",
		"detail": ""
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusForbidden, w.Code)

}

func TestHandleTaskCreateBadContentInt(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	requestBody := []byte(`{"content": 43}`)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	srv.handleTaskCreate()(w, req)

	expectedResp := `{
		"error": "Cannot decode task body from json",
		"detail": "json: cannot unmarshal number into Go struct field request.content of type string"
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)

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

func TestHandleTaskEditBodyWrongType(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	taskID := "12"

	requestBody := []byte(`{"content": 40}`)
	req := httptest.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": taskID})
	w := httptest.NewRecorder()
	srv.handleTaskEdit()(w, req)

	expectedResp := `{
		"error": "Cannot parse task body",
		"detail": "json: cannot unmarshal number into Go struct field request.content of type string"
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleTaskEditContentEmpty(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	taskID := "12"

	requestBody := []byte(`{"content": ""}`)
	req := httptest.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": taskID})
	w := httptest.NewRecorder()
	srv.handleTaskEdit()(w, req)

	expectedResp := `{
		"error": "Key 'content' cannot be empty",
		"detail": ""
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusForbidden, w.Code)
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

func TestHandleChangeTaskStateBadID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()
	srv := &server{
		DB: &database.DBStore{DB: db},
	}

	taskID := 12

	query := "SELECT id, content, state FROM tasks WHERE id = \\$1"
	mock.ExpectQuery(query).WithArgs(taskID).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("PUT", "/tasks/state/12", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "12"})
	w := httptest.NewRecorder()
	srv.handleTaskState()(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
	expectedResp := `{
		"error": "Task not found",
		"detail": "sql: no rows in result set"
	  }`
	assert.JSONEq(t, expectedResp, w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}
