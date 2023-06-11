package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wilianto/bababos-pricing-engine/customer"
	"github.com/wilianto/bababos-pricing-engine/price"
	"github.com/wilianto/bababos-pricing-engine/transport"
)

type PriceIntegrationTestSuite struct {
	suite.Suite
	DB *sqlx.DB
}

func TestPriceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PriceIntegrationTestSuite))
}

func (s *PriceIntegrationTestSuite) SetupSuite() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:password@localhost:5432/pricing_engine?sslmode=disable")
	require.NoError(s.T(), err)
	s.DB = db
}

func (s *PriceIntegrationTestSuite) TearDownSuite() {
	_, err := s.DB.Exec("DELETE FROM price")
	require.NoError(s.T(), err)

	_, err = s.DB.Exec("DELETE FROM item")
	require.NoError(s.T(), err)

	_, err = s.DB.Exec("DELETE FROM supplier")
	require.NoError(s.T(), err)

	_, err = s.DB.Exec("DELETE FROM customer")
	require.NoError(s.T(), err)
}

func (s *PriceIntegrationTestSuite) TestGetPrice() {
	skuID := "1234568"
	customerID, supplierID := s.seedData(skuID)

	path := fmt.Sprintf("/price?customer_id=%v&sku_id=%v&qty=120", customerID, skuID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, rec)
	handler := &transport.HttpHandler{
		BasicPricing: &price.BasicPricingStrategy{
			PriceRepository:    &price.PriceRepository{DB: s.DB},
			CustomerRepository: &customer.CustomerRepository{DB: s.DB},
		},
	}

	err := handler.GetPrice(ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, rec.Code)
	require.Equal(
		s.T(),
		fmt.Sprintf("{\"suggested_price\":1100,\"recommended_supplier_id\":%v}\n", supplierID),
		rec.Body.String(),
	)
}

func (s *PriceIntegrationTestSuite) seedData(skuID string) (int64, int64) {
	custRow, err := s.DB.NamedQuery("INSERT INTO customer (address, city, state) VALUES (:address, :city, :state) RETURNING id", map[string]interface{}{
		"address": "Holis",
		"city":    "bandung",
		"state":   "singapore",
	})
	require.NoError(s.T(), err)
	var customerID int64
	if custRow.Next() {
		custRow.Scan(&customerID)
	}
	require.NoError(s.T(), err)

	_, err = s.DB.NamedExec("INSERT INTO item (sku_id, name) VALUES (:sku_id, :name)", map[string]interface{}{
		"sku_id": skuID,
		"name":   "besi",
	})
	require.NoError(s.T(), err)

	supplierRow, err := s.DB.NamedQuery("INSERT INTO supplier (address, city, state) VALUES (:address, :city, :state) RETURNING id", map[string]interface{}{
		"address": "TKI",
		"city":    "bandung",
		"state":   "singapore",
	})
	var supplierID int64
	if supplierRow.Next() {
		supplierRow.Scan(&supplierID)
	}
	require.NoError(s.T(), err)

	_, err = s.DB.NamedExec("INSERT INTO price (sku_id, supplier_id, price) VALUES (:sku_id, :supplier_id, :price)", map[string]interface{}{
		"sku_id":      skuID,
		"supplier_id": supplierID,
		"price":       1000,
	})
	require.NoError(s.T(), err)

	return customerID, supplierID
}
