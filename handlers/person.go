package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lincentpega/kodalib-rpc-db/data"
	"github.com/lincentpega/kodalib-rpc-db/utils"
)

type Person struct {
	db *sql.DB
	l  *log.Logger
}

func NewPerson(db *sql.DB, l *log.Logger) *Person {
	return &Person{db, l}
}

func (p *Person) GetPersonByName(rw http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&rw)
	vars := mux.Vars(r)
	name := vars["name"]
	p.l.Printf("Processing GET reqiest on /api/persons/" + name)
	persons, err := data.GetPersons(p.db, name)
	if err != nil {
		http.Error(rw, "Unable to find person with name: "+name, http.StatusNotFound)
		return
	}
	if err = persons.ToJSON(rw); err != nil {
		http.Error(rw, "Unable to return person", http.StatusInternalServerError)
	}
}
