package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
