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

func (store *DBStore) CreateTask(t *Task) (int64, error) {
	var id int64
	err := store.db.QueryRow("INSERT INTO tasks (content,state) VALUES ($1, $2) RETURNING id", t.Content, t.State).Scan(&id)
	return id, err
}

func (store *DBStore) DeleteTask(taskID int) error {
	// Exécuter la requête DELETE
	_, err := store.db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	return err
}

func (store *DBStore) EditTask(taskID int, content string) error {
	// Check if the row with the specified ID exists
	var exists bool
	err := store.db.QueryRow("SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)", taskID).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check if row exists. err=%v\n", err)
		return err
	}
	if !exists {
		// ID not found, return a custom error
		err = &CustomError{
			Message: fmt.Sprintf("Row with ID %d not found", taskID),
		}
		log.Printf("Row with ID %v does not exist", taskID)
		return err
	}
	_, err = store.db.Exec("UPDATE tasks SET content = $1 WHERE id = $2", content, taskID)
	return err
}
