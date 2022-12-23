package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/lincentpega/kodalib-rpc-db/handlers"
)

type App struct {
	fh     *handlers.Films
	ph     *handlers.Person
	dbh    *handlers.DB
	Router *mux.Router
	DB     *sql.DB
	l      *log.Logger
}

func (a *App) Initialize(db_host, db_user, db_name, db_pass, address string, l *log.Logger) {
	cs := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", db_host, db_user, db_name, db_pass)
	db, err := sql.Open("postgres", cs)
	if err != nil {
		log.Fatal(err)
	}

	a.DB = db
	a.l = l
	a.dbh = handlers.NewDB(db, l)
	a.fh = handlers.NewFilms(db, l)
	a.ph = handlers.NewPerson(db, l)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) SetRouter(r *mux.Router) {
	a.Router = r
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/films/{title}", a.fh.GetFilmsByTitle).Methods(http.MethodGet)
	a.Router.HandleFunc("/api/persons/{name}", a.ph.GetPersonByName).Methods(http.MethodGet)
	a.Router.HandleFunc("/api/films/{title}", a.fh.DeleteFilmsByTitle).Methods(http.MethodDelete)
	a.Router.HandleFunc("/api/persons/{name}", a.ph.DeletePersonsByName).Methods(http.MethodDelete)
	a.Router.HandleFunc("/api/create_db/{title}", a.dbh.CreateDB).Methods(http.MethodPost)
	a.Router.HandleFunc("/api/delete_db/{title}", a.dbh.DeleteDB).Methods(http.MethodPost)
	a.Router.HandleFunc("/api/table/{title}", a.dbh.TruncateTable).Methods(http.MethodDelete)
	a.Router.HandleFunc("/api/table/all", a.dbh.TruncateAll).Methods(http.MethodDelete)
}

func (a *App) Run(address string, l *log.Logger) {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	s := &http.Server{
		Addr:         address,
		Handler:      a.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     l,
	}

	go func() {
		l.Println("Starting server on port " + os.Getenv("PORT"))
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	s.Shutdown(ctx)
	l.Println("shutting down")
	os.Exit(0)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	l := log.New(os.Stdout, "koda-db ", log.LstdFlags)

	a := &App{}
	a.Initialize(
		os.Getenv("DBHOST"),
		os.Getenv("DBUSER"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPASS"),
		os.Getenv("PORT"),
		l,
	)

	a.Run(":"+os.Getenv("PORT"), l)
}
