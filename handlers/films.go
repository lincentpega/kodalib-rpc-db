package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lincentpega/kodalib-rpc-db/data"
)

type Films struct {
	db *sql.DB
	l  *log.Logger
}

func NewFilms(db *sql.DB, l *log.Logger) *Films {
	return &Films{db, l}
}

func (f *Films) GetFilmByTitle(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	f.l.Printf("%v", vars)
	title := vars["title"]
	film, err := data.GetFilm(f.db, title)
	if err != nil {
		http.Error(rw, "unable to find film with title "+title, http.StatusNotFound)
		return
	}
	if err = film.ToJSON(rw); err != nil {
		http.Error(rw, "unable to return film", http.StatusInternalServerError)
	}
}
