package database_test

import (
	"fmt"
	"reflect"
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
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(1, "Task 1", false).
		AddRow(2, "Task 2", false)

	mock.ExpectQuery("SELECT id, content, state FROM tasks").WillReturnRows(rows)

	tasks, err := srv.DB.GetTaskList()
	if err != nil {
		t.Fatalf("Error while executing GetTaskList : %s", err)
	}

	expectedTasks := []*database.Task{
		{ID: 1, Content: "Task 1", State: false},
		{ID: 2, Content: "Task 2", State: false},
	}

	assert.Equal(t, expectedTasks, tasks, "Tasks does not correspond")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Different behavior between expectation and result : %s", err)
	}
}

func TestGetTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "content", "state"}).
		AddRow(1, "Task 1", false).
		AddRow(2, "Task 2", false)

	query := "SELECT id, content, state FROM tasks WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

	task, err := srv.DB.GetTask(1)
	if err != nil {
		t.Fatalf("Error while executing GetTask : %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}

	expectedTask := database.Task{ID: 1, Content: "Task 1", State: false}
	assert.Equal(t, &expectedTask, task, "Task does not correspond")
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()

	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	task := &database.Task{
		ID:      0,
		Content: "test task",
		State:   false,
	}

	insert := "INSERT INTO tasks (content,state) VALUES ($1, $2) RETURNING id"
	mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(task.Content, task.State).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := srv.DB.CreateTask(task)
	if err != nil {
		t.Fatalf("Error while creating task : %s", err)
	}

	if id != 1 {
		t.Fatalf("Bad task ID, wanted 1 but got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()
	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	taskID := 12
	delete := "DELETE FROM tasks WHERE id = \\$1"
	mock.ExpectExec(delete).WithArgs(taskID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = srv.DB.DeleteTask(taskID)
	if err != nil {
		t.Errorf("Error while deleting task : %v", err)
	}
}

func TestDeleteTaskBadID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()
	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	taskID := 12
	delete := "DELETE FROM tasks WHERE id = \\$1"
	mock.ExpectExec(delete).WithArgs(taskID).WillReturnResult(sqlmock.NewResult(0, 0))

	err = srv.DB.DeleteTask(taskID)
	expectedError := fmt.Sprintf("task with ID %d does not exist", taskID)
	if err.Error() != expectedError {
		t.Fatalf("Function DeleteTask returned bad error message")
	}
}

func TestEditTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()
	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

	taskID := 123
	content := "task content"
	query := "SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(taskID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectExec("UPDATE tasks SET content = \\$1 WHERE id = \\$2").
		WithArgs(content, taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = srv.DB.EditTask(taskID, content)
	if err != nil {
		t.Fatalf("Error while editing task : %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
}

func TestChangeTaskState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while creating mock : %s", err)
	}
	defer db.Close()
	srv := router.NewServer()
	srv.DB = &database.DBStore{DB: db}

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

	expectedTask := &database.Task{
		ID:      int64(taskID),
		Content: "Task 1",
		State:   state,
	}

	task, err := srv.DB.ChangeTaskState(taskID)
	if err != nil {
		t.Fatalf("Error while changing task state : %v", err)
	}
	if !reflect.DeepEqual(task, expectedTask) {
		t.Fatalf("Wrong task, got %v, wanted %v", task, expectedTask)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met : %s", err)
	}
}
