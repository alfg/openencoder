package data

import (
	_ "database/sql" // Database.
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver.
)

const dialect = "postgres"

// var connectionString = os.Getenv("DATABASE")
var connectionString = ""

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "enc"
)

func init() {
	// if connectionString == "" {
	// 	// connectionString = settings.Get().Database.Database
	// }
	connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// ConnectDB Connects to sqlite3 database
func ConnectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect(dialect, connectionString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to DB.")
	return db, err
}