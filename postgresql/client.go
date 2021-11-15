package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func GetPostgresClient() *sqlx.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		"user",
		"mypwd",
		"localhost",
		"calendar",
	)

	db, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		log.Fatal("error connection to Postgres : ", err.Error())
	}

	return db
}
