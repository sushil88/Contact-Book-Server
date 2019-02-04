package db

import (
	"log"

	"contactBook/conf"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Load pq (postgres driver) only when this db package has been loaded
)

var DB *sqlx.DB // DB used for database connection handle

// Get returns either an already loaded gorm database connection or sets one up
func Get() *sqlx.DB {
	var err error

	if DB == nil || DB.Ping() != nil {
		DB, err = sqlx.Connect("postgres", conf.Get().Database.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}
	}

	return DB
}
