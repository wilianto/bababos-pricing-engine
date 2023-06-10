package price

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PriceRepositoryTestSuite struct {
	suite.Suite
	DB *sqlx.DB
}

func TestPriceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PriceRepositoryTestSuite))
}

func (suite *PriceRepositoryTestSuite) SetupSuite() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:password@localhost:5432/pricing_engine?sslmode=disable")
	require.NoError(suite.T(), err)
	suite.DB = db
}

func (s *PriceRepositoryTestSuite) TearDownSuite() {
	_, err := s.DB.Exec("DELETE FROM price")
	require.NoError(s.T(), err)

	_, err = s.DB.Exec("DELETE FROM item")
	require.NoError(s.T(), err)

	_, err = s.DB.Exec("DELETE FROM supplier")
	require.NoError(s.T(), err)
}

func (s *PriceRepositoryTestSuite) TestGetCustomer() {
	singapore := "singapore"
	skuID := "1234568"
	_, err := s.DB.NamedExec("INSERT INTO item (sku_id, name) VALUES (:sku_id, :name)", map[string]interface{}{
		"sku_id": skuID,
		"name":   "besi",
	})
	require.NoError(s.T(), err)

	supplierID1 := s.InsertSupplier("supplier 1", singapore)
	supplierID2 := s.InsertSupplier("supplier 2", singapore)

	s.InsertPrice(skuID, supplierID1, 1000)
	s.InsertPrice(skuID, supplierID2, 1500)

	tests := []struct {
		name    string
		skuID   string
		state   string
		want    *RecommendedPrice
		wantErr error
	}{
		{
			name:  "when sku id and state is found return lowest price",
			skuID: skuID,
			state: singapore,
			want:  &RecommendedPrice{Price: 1000, SupplierID: supplierID1},
		},
		{
			name:    "when sku id and state is not found return error",
			skuID:   skuID,
			state:   "jakarta",
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &PriceRepository{
				DB: s.DB,
			}
			got, err := r.FindLowestPrice(tt.skuID, tt.state)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func (s *PriceRepositoryTestSuite) InsertSupplier(address string, state string) int64 {
	row, err := s.DB.NamedQuery("INSERT INTO supplier (address, city, state) VALUES (:address, :city, :state) RETURNING id", map[string]interface{}{
		"address": address,
		"city":    "bandung",
		"state":   state,
	})
	require.NoError(s.T(), err)

	var supplierID int64
	if row.Next() {
		row.Scan(&supplierID)
	}

	return supplierID
}

func (s *PriceRepositoryTestSuite) InsertPrice(skuID string, supplierID int64, price float64) {
	_, err := s.DB.NamedExec("INSERT INTO price (sku_id, supplier_id, price) VALUES (:sku_id, :supplier_id, :price)", map[string]interface{}{
		"sku_id":      skuID,
		"supplier_id": supplierID,
		"price":       price,
	})
	require.NoError(s.T(), err)
}
