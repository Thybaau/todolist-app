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
