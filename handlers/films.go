package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lincentpega/kodalib-rpc-db/data"
	"github.com/lincentpega/kodalib-rpc-db/utils"
)

type Films struct {
	db *sql.DB
	l  *log.Logger
}

func NewFilms(db *sql.DB, l *log.Logger) *Films {
	return &Films{db, l}
}

func (f *Films) GetFilms(rw http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&rw)

}

func (f *Films) GetFilmsByTitle(rw http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&rw)
	vars := mux.Vars(r)
	title := vars["title"]
	f.l.Printf("Processing GET reqiest on /api/films/" + title)
	films, err := data.GetFilmsByTitle(f.db, title)
	if err != nil {
		http.Error(rw, "Unable to find film with title: "+title, http.StatusNotFound)
		return
	}
	if err = films.ToJSON(rw); err != nil {
		http.Error(rw, "Unable to return film", http.StatusInternalServerError)
	}
}

func (f *Films) DeleteFilmsByTitle(rw http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&rw)
	vars := mux.Vars(r)
	title := vars["title"]
	f.l.Printf("Processing DELETE request on /api/films/" + title)
	res := data.DeleteFilms(f.db, f.l, title)
	if res {
		return
	} else {
		http.Error(rw, "There is no films with title: "+title, http.StatusNotFound)
		return
	}
}
