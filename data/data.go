package data

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/lib/pq"
)

//Db pointer on databese object
var Db *sql.DB


const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "docsampler"
)



func open() {
	var err error
	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	
	Db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

func close() {
	Db.Close()
}

