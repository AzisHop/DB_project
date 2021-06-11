package main

import (
	"DB_project/db"
	"fmt"
	"log"
)

func main() {
	postgres, err := db.NewDb("user=postgres dbname=postgres password=admin host=127.0.0.1 port=5432 sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
}