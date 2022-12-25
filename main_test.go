package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var a App

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PASS"),
		os.Getenv("PORT"),
		log.New(os.Stdout, "test-kodalib-api", log.LstdFlags),
	)

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)

}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM films")
	a.DB.Exec("ALTER SEQUENCE films_id_seq RESTART WITH 1")
}

const tableCreationQuery = `create table if not exists films
(
    id               serial,
    imdb_id          text,
    title            text not null,
    poster           text,
    year             smallint,
    duration         text,
    plot             text,
    imdb_rating      text,
    kodalib_rating   text,
    budget           text,
    gross_worldwide  text,
    youtube_trailer  text,
    "ThumbnailUrl"   text,
    count_of_changes integer
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest(http.MethodGet, "/api/films", nil)
	resp := executeRequest(req)
	checkResponseCode(t, http.StatusOK, resp.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentFilm(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest(http.MethodGet, "/api/films/111", nil)
	resp := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, resp.Code)

	if resp.Body.String() != "Unable to find film with title: 111\n" {
		t.Errorf("Expected the value of the response to 'Unable to find film with title: 111. Got '%s'", resp.Body.String())
	}
}
