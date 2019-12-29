package data

import (
	_ "database/sql" // Database.
	"fmt"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/logging"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // Postgres driver.
)

const (
	driverName = "postgres"
)

var (
	connectionString string
	conn             *sqlx.DB
	log              = logging.Log
	connectionFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
)

// ConnectDB Connects to postgres database
func ConnectDB() (*sqlx.DB, error) {
	var err error
	if connectionString == "" {
		log.Info("setting database connectionString.")
		var (
			host     = config.Get().DatabaseHost
			port     = config.Get().DatabasePort
			user     = config.Get().DatabaseUser
			password = config.Get().DatabasePassword
			dbname   = config.Get().DatabaseName
		)
		connectionString = fmt.Sprintf(connectionFormat, host, port, user, password, dbname)
	}

	if conn, err = sqlx.Connect(driverName, connectionString); err != nil {
		log.Panic(err)
	}
	return conn, err
}
