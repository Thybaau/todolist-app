package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBStore struct {
	db *sql.DB
}

// Tasks structs
type Task struct {
	ID      int64  `db:"id"`
	Content string `db:"content"`
	State   bool   `db:"state"`
}

type CustomError struct {
	Message string
}

// Implémentation de l'interface error pour CustomError
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error : %s", e.Message)
}

func (store *DBStore) Connect(host string, port int, user, password, dbname string) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	// Ping to check connection
	err = db.Ping()
	if err != nil {
		return err
	}
	log.Printf("Connected to Postgre DB %s", dbname)
	store.db = db
	return nil
}

func (store *DBStore) Close() error {
	return store.db.Close()
}

func (store *DBStore) GetTaskList() ([]*Task, error) {
	rows, err := store.db.Query("SELECT id, content, state FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Add values to each Task structs and return list of structs
	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Content, &t.State); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil

}
func (store *DBStore) GetTask(id int) (*Task, error) {
	row := store.db.QueryRow("SELECT id, content, state FROM tasks WHERE id = $1", id)

	var task Task
	if err := row.Scan(&task.ID, &task.Content, &task.State); err != nil {
		return nil, err
	}

	return &task, nil
}

func (store *DBStore) CreateTask(t *Task) (int64, error) {
	var id int64
	err := store.db.QueryRow("INSERT INTO tasks (content,state) VALUES ($1, $2) RETURNING id", t.Content, t.State).Scan(&id)
	return id, err
}

func (store *DBStore) DeleteTask(taskID int) error {
	result, err := store.db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	//Check number of lines affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no lines affected, it means ID didn't exist
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d does not exist", taskID)
	}
	return nil
}

func (store *DBStore) EditTask(taskID int, content string) error {
	// Check if the row with the specified ID exists
	var exists bool
	err := store.db.QueryRow("SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)", taskID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		// ID not found, return a custom error
		err = &CustomError{
			Message: fmt.Sprintf("row with ID %d not found", taskID),
		}
		return err
	}
	_, err = store.db.Exec("UPDATE tasks SET content = $1 WHERE id = $2", content, taskID)
	return err
}

func (store *DBStore) ChangeTaskState(taskID int) (*Task, error) {
	task, err := store.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	newState := !task.State
	_, err = store.db.Exec("UPDATE tasks SET state = $1 WHERE id = $2", newState, taskID)
	if err != nil {
		return nil, err
	}
	task, err = store.GetTask(taskID)
	return task, err
}
