package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	runMigration()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", healthCheck)

	e.Logger.Fatal(e.Start(":8080"))
}

func healthCheck(e echo.Context) error {
	return e.JSON(200, map[string]string{"status": "ok"})
}

func runMigration() {
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
