package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Person struct {
	Id           int    `json:"id"`
	PersonImdbId string `json:"person_imdb_id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Summary      string `json:"summary"`
	BirthDate    string `json:"birth_date"`
	DeathDate    string `json:"death_date"`
	Height       string `json:"height"`
}

func (f *Person) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(f)
}

func (f *Person) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

type Persons []Person

func (f *Persons) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

func (f *Persons) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(f)
}

func GetPersons(db *sql.DB, l *log.Logger, name string) (Persons, error) {
	q := "SELECT id, person_imdb_id, name, image, summary, birth_date, death_date, height FROM get_persons_by_name($1);"

	rows, err := db.Query(q, name)
	if err != nil {
		l.Fatal(err)
	}
	defer rows.Close()

	var ps Persons

	isAnyMore := rows.Next()
	if !isAnyMore {
		return nil, fmt.Errorf("unable to find person with name: %s", name)
	}

	for isAnyMore {
		var p Person
		rows.Scan(&p.Id, &p.PersonImdbId, &p.Name, &p.Image, &p.Summary, &p.BirthDate, &p.DeathDate, &p.Height)

		// bacause of non-void stored function call, which returns sth with id 0, if person not found
		if p.Id == 0 {
			return nil, fmt.Errorf("unable to find person with name : %s", name)
		}

		ps = append(ps, p)

		isAnyMore = rows.Next()
	}

	return ps, nil
}

func DeletePersons(db *sql.DB, l *log.Logger, name string) bool {
	q := "SELECT * FROM delete_persons_by_name_if_exists($1)"

	rows, err := db.Query(q, name)
	if err != nil {
		l.Fatal(err)
	}
	defer rows.Close()

	var res int

	if !rows.Next() {
		return false
	}

	rows.Scan(&res)
	if res == 0 {
		return true
	} else {
		return false
	}

}
