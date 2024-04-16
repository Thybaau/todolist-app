package router

import (
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

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestHandleTaskListOneTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(1, "Task 1", false).
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

	assert.Equal(t, http.StatusOK, w.Code)
}

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

	assert.Equal(t, http.StatusOK, w.Code)
}
