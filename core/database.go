package core

import (
	"database/sql"

	"github.com/anand-kay/linkedin-clone/config"

	// Postgres driver
	_ "github.com/lib/pq"
)

// db - Database instance
var db *sql.DB

// newDB - Opens a new database connection
func newDB(config *config.Config) (*sql.DB, error) {
	connStr := "host=" + config.DbHost + " port=" + config.DbPort + " user=" + config.DbUser + " password=" + config.DbPassword + " dbname=" + config.DbName + " sslmode=" + config.DbSslmode

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
