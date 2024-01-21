package tests

import (
	"regexp"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Thybaau/todolist-app/database"
	"github.com/Thybaau/todolist-app/router"
)

func TestGetTaskList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Error while creating database mock")
	defer db.Close()

	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(1, "Task 1", false).
		AddRow(2, "Task 2", false)

	mock.ExpectQuery("SELECT id, content, state FROM tasks").WillReturnRows(rows)

	tasks, err := srv.DB.GetTaskList()
	assert.NoError(t, err, "Error while executing GetTaskList")

	expectedTasks := []*database.Task{
		{ID: 1, Content: "Task 1", State: false},
		{ID: 2, Content: "Task 2", State: false},
	}

	assert.Equal(t, expectedTasks, tasks, "Tasks does not correspond")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Different behavior between expectation and result")
}

func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Error while creating database mock")
	defer db.Close()

	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(1, "Task 1", false).
		AddRow(2, "Task 2", false)

	query := "SELECT id, content, state FROM tasks WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

	task, err := srv.DB.GetTask(1)
	assert.NoError(t, err, "Error while executing GetTask")

	expectedTask := database.Task{ID: 1, Content: "Task 1", State: false}
	assert.Equal(t, &expectedTask, task, "Task does not correspond")
}
