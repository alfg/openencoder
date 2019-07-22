package data

import (
	_ "database/sql" // Database.
	"fmt"
	"log"

	"github.com/alfg/openencoder/api/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // Postgres driver.
	// _ "github.com/mattn/go-sqlite3"
)

var (
	connectionString = ""
	conn             *sqlx.DB
)

// ConnectDB Connects to postgres database
func ConnectDB() (*sqlx.DB, error) {
	var err error
	if connectionString == "" {
		fmt.Println("connection not set. setting now.")
		var (
			host     = config.Get().DatabaseHost
			port     = config.Get().DatabasePort
			user     = config.Get().DatabaseUser
			password = config.Get().DatabasePassword
			dbname   = config.Get().DatabaseName
		)
		connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	if conn, err = sqlx.Connect("postgres", connectionString); err != nil {
		log.Panic(err)
	}
	return conn, err
}
