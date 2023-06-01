package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// TODO: read from env
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/pricing_engine?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("error opening db conn due to %s", err.Error()))
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("error creating postgres instance due to %s", err.Error()))
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migration", "postgres", driver)
	if err != nil {
		panic(fmt.Sprintf("error creating migration instance due to %s", err.Error()))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("error running migration due to %s", err.Error()))
	}

	fmt.Println("running migration success")
}
