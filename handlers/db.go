package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DB struct {
	db *sql.DB
	l  *log.Logger
}

func NewDB(db *sql.DB, l *log.Logger) *DB {
	return &DB{db, l}
}

func (db *DB) CreateDB(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	q := "SELECT * FROM create_base($1)"
	if _, err := db.db.Query(q, title); err != nil {
		db.l.Fatal(err)
		http.Error(rw, "Unable to create database with title: "+title, http.StatusInternalServerError)
		return
	}

}

func (db *DB) DeleteDB(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	q := "CALL drop_base($1)"
	if _, err := db.db.Query(q, title); err != nil {
		db.l.Fatal(err)
		http.Error(rw, "Unable to delete database with title: "+title, http.StatusInternalServerError)
		return
	}
}

func (db *DB) TruncateTable(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	q := "SELECT * FROM truncate_if_exists($1)"

	if _, err := db.db.Query(q, title); err != nil {
		db.l.Fatal(err)
		http.Error(rw, "Unable to truncate table: "+title, http.StatusInternalServerError)
		return
	}
}

func (db *DB) TruncateAll(rw http.ResponseWriter, r *http.Request) {
	q := "SELECT * FROM truncate_all()"

	if _, err := db.db.Query(q); err != nil {
		db.l.Fatal(err)
		http.Error(rw, "Unable to truncate tables", http.StatusInternalServerError)
		return
	}
}
