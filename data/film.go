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

func GetFilms(db *sql.DB, title string) (Films, error) {
	q := "SELECT * FROM get_films_by_title($1);"

	rows, err := db.Query(q, title)
	if err != nil {
		log.Fatal(err)
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

		// bacause of non-void stored function call, which returns sth with id 0, if film not found
		if f.Id == 0 {
			return nil, fmt.Errorf("unable to find film with title: %s", title)
		}

		fs = append(fs, f)
		isAnyMore = rows.Next()
	}

	return fs, nil
}
