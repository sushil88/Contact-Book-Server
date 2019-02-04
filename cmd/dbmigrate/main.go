package main

import (
	"flag"
	"fmt"

	"github.com/DavidHuie/gomigrate"

	"contactBook/db"
)

var method = "up"

func init() {
	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		method = arg
	}
}

func main() {
	db := db.Get().DB

	migrator, err := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
	if err != nil {
		panic(err)
	}
	if method == "rollback" {
		err = migrator.Rollback()
	} else if method == "up" {
		err = migrator.Migrate()
	} else {
		panic("The migrate method of " + method + ` is invalid.  "up" and "rollback" are allowed`)
	}

	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

}
