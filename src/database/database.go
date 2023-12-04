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

type Task struct {
	Content string `db:"content"`
	State   bool   `db:"state"`
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

func (store *DBStore) CreateTask(t *Task) (int64, error) {
	var id int64
	err := store.db.QueryRow("INSERT INTO tasks (content,state) VALUES ($1, $2) RETURNING id", t.Content, t.State).Scan(&id)
	return id, err
}