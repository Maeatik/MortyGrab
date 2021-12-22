package models

import (
	"database/sql"
	"fmt"
)

const (
	host     = 	"localhost"
	port     = 	5432
	user     = 	"postgres"
	password = 	"pass"
	dbname   = 	"MortyGRAB"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//psqlInfo := "host=localhost port=5432 user=postgres password=pass dbname=MortyGRAB sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}