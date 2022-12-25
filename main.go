package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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
	p := new([3]int)
	fmt.Printf("%v", *p)
}
