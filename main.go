package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/wilianto/bababos-pricing-engine/customer"
	"github.com/wilianto/bababos-pricing-engine/price"
	"github.com/wilianto/bababos-pricing-engine/transport"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runServer()
		return
	}

	switch args[0] {
	case "db:seed":
		runSeed()
	case "db:migrate":
		runMigration()
	default:
		runServer()
	}
}

func runServer() {
	handler := initHandler()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", healthCheck)
	e.GET("/price", handler.GetPrice)

	e.Logger.Fatal(e.Start(":8080"))
}

func initHandler() transport.HttpHandler {
	db, err := sqlx.Connect("postgres", "postgres://postgres:password@localhost:5432/pricing_engine?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("error connecting to db due to %s", err.Error()))
	}

	customerRepo := customer.CustomerRepository{DB: db}
	priceRepo := price.PriceRepository{DB: db}
	basicPricing := price.BasicPricingStrategy{PriceRepository: &priceRepo, CustomerRepository: &customerRepo}
	surgePricing := price.SurgePricingStrategy{PriceRepository: &priceRepo, CustomerRepository: &customerRepo}
	handler := transport.HttpHandler{BasicPricing: &basicPricing, SurgePricing: &surgePricing}
	return handler
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

func runSeed() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:password@localhost:5432/pricing_engine?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("error connecting to db due to %s", err.Error()))
	}

	custRow, err := db.Query("INSERT INTO customer (address, city, state) VALUES ('Jl. Raya Bogor', 'Bogor', 'Jawa Barat') RETURNING id")
	if err != nil {
		fmt.Printf("error inserting customer %v\n", err)
	}
	var custID int64
	if custRow.Next() {
		custRow.Scan(&custID)
	}

	supplierRow, err := db.Query("INSERT INTO supplier (address, city, state) VALUES ('Jl. Raya Bogor', 'Bogor', 'Jawa Barat') RETURNING id")
	if err != nil {
		fmt.Printf("error inserting supplier %v\n", err)
	}
	var supplierID int64
	if supplierRow.Next() {
		supplierRow.Scan(&supplierID)
	}

	skuID := "12345"
	_, err = db.NamedQuery("INSERT INTO item (sku_id, name) VALUES (:sku_id, 'Besi')", map[string]interface{}{"sku_id": skuID})
	if err != nil {
		fmt.Printf("error inserting item %v\n", err)
	}

	_, err = db.NamedQuery("INSERT INTO price (sku_id, supplier_id, price) VALUES (:sku_id, :supplier_id, :price)", map[string]interface{}{
		"sku_id":      skuID,
		"supplier_id": supplierID,
		"price":       10000,
	})
	if err != nil {
		fmt.Printf("error inserting price %v\n", err)
	}

	fmt.Println("============ running seed success ============")
	fmt.Printf("customer id: %d\n", custID)
	fmt.Printf("supplier id: %d\n", supplierID)
	fmt.Printf("sku id: %s\n", skuID)
	fmt.Println("============ running seed success ============")
}
