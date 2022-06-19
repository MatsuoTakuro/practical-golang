package ch9

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB
var dbErr error

func init() {
	db, dbErr = sql.Open("pgx", "host=localhost port=5432 user=testuser dbname=testdb password=pass sslmode=disable")
	if nil != dbErr {
		log.Fatal(dbErr)
	}

	cmdU := `CREATE TABLE IF NOT EXISTS users(
		user_id varchar(32) NOT NULL,
		user_name varchar(100) NOT NULL,
		created_at timestamp with time zone,
		CONSTRAINT pk_users PRIMARY KEY (user_id)
	)`

	_, err := db.Exec(cmdU)
	if err != nil {
		log.Fatal(err)
	}

	InsertInitialUsers()
}

func InsertInitialUsers() {
	cmdU := `INSERT INTO users(
		user_id,
		user_name,
		created_at
	)
	VALUES('0001', 'Gopher', '2020-07-10 00:00:00.000000+00'),
				('0002', 'Ferris', '2020-07-11 00:00:00.000000+00')`

	// TODO: handle this -> ERROR: duplicate key value violates unique constraint "pk_users" (SQLSTATE 23505)
	_, err := db.Exec(cmdU)
	if err != nil {
		log.Println(err)
	}
}

func Sub() {
	queryMultiLines()
	querySingleLine()
}
