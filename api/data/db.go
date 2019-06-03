package data

import (
	_ "database/sql" // Database.
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	// _ "github.com/lib/pq" // Postgres driver.
	_ "github.com/mattn/go-sqlite3"
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
	// connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
}

var schema = `
DROP TABLE IF EXISTS user;
CREATE TABLE user (
	id    INTEGER PRIMARY KEY,
);
`

// ConnectDB Connects to sqlite3 database
func ConnectDB() (*sqlx.DB, error) {
	// db, err := sqlx.Connect(dialect, connectionString)
	db, err := sqlx.Connect("sqlite3", "enc.db")

	// db.MustExec(schema)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to DB.")
	return db, err
}
