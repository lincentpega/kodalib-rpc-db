package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Film struct {
	Id             int    `json:"id"`
	ImdbId         string `json:"imdb_id"`
	Title          string `json:"title"`
	Poster         string `json:"poster"`
	Year           int    `json:"year"`
	Duration       string `json:"duration"`
	Plot           string `json:"plot"`
	ImdbRating     string `json:"imdb_rating"`
	KodalibRating  string `json:"kodalib_rating"`
	Budget         string `json:"budget"`
	GrossWorldwide string `json:"gross_worldwide"`
	YoutubeTrailer string `json:"youtube_trailer"`
	ThumbnailUrl   string `json:"thumbnail_url"`
}

func (f *Film) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(f)
}

func (f *Film) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

type Films []Film

func (f *Films) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

func (f *Films) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(f)
}

func GetAllFilms(db *sql.DB) (Films, error) {
	q := `SELECT * FROM 
	id, imdb_id, title, poster, year, duration, plot, 
	imdb_rating, kodalib_rating, budget, gross_worldwide, youtube_trailer, "ThumbnailUrl"`

	rows, err := db.Query(q)
	if err != nil {
		err = fmt.Errorf("unable to retreive films, got error: %v", err)
		return nil, err
	}

	var films Films

	for rows.Next() {
		var f Film
		rows.Scan(&f.Id, &f.ImdbId, &f.Title, &f.Poster, &f.Year, &f.Duration,
			&f.Plot, &f.ImdbRating, &f.KodalibRating, &f.Budget, &f.GrossWorldwide,
			&f.YoutubeTrailer, &f.ThumbnailUrl)
		films = append(films, f)
	}

	return films, nil
}

func GetFilmsByTitle(db *sql.DB, title string) (Films, error) {
	q := `SELECT id, imdb_id, title, poster, year, duration, plot, imdb_rating, kodalib_rating, budget, gross_worldwide, youtube_trailer, "ThumbnailUrl" 
	FROM films 
	WHERE films.title = $1`

	rows, err := db.Query(q, title)
	if err != nil {
		err = fmt.Errorf("unable to retreive films, got error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var fs Films

	isAnyMore := rows.Next()
	if !isAnyMore {
		return nil, fmt.Errorf("unable to find film with title: %s", title)
	}
	for isAnyMore {
		var f Film
		rows.Scan(&f.Id, &f.ImdbId, &f.Title, &f.Poster, &f.Year, &f.Duration,
			&f.Plot, &f.ImdbRating, &f.KodalibRating, &f.Budget, &f.GrossWorldwide,
			&f.YoutubeTrailer, &f.ThumbnailUrl)

		fs = append(fs, f)
		isAnyMore = rows.Next()
	}

	return fs, nil
}

func DeleteFilms(db *sql.DB, l *log.Logger, title string) bool {
	q := "SELECT * FROM delete_films_by_title_if_exists($1)"

	rows, err := db.Query(q, title)
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
