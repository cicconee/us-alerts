package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var host = "0.0.0.0"
var port = "5432"
var user = "usalert"
var password = "password"
var dbname = "usalertsdb"

func main() {
	fmt.Println("US Alerts - Get national weather alerts.")

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v\n", err)
	}
}
