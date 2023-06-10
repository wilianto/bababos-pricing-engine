package customer

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	DB *sqlx.DB
}

func TestCustomerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}

func (suite *CustomerRepositoryTestSuite) SetupSuite() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:password@localhost:5432/pricing_engine?sslmode=disable")
	require.NoError(suite.T(), err)
	suite.DB = db
}

func (s *CustomerRepositoryTestSuite) TearDownSuite() {
	_, err := s.DB.Exec("DELETE FROM customer")
	require.NoError(s.T(), err)
}

func (s *CustomerRepositoryTestSuite) TestGetCustomer() {
	row, err := s.DB.Query("INSERT INTO customer (address, city, state) VALUES ('address', 'city', 'state') RETURNING id", Customer{})
	require.NoError(s.T(), err)

	var insertedCustomerID int64
	if row.Next() {
		row.Scan(&insertedCustomerID)
	}

	tests := []struct {
		name       string
		customerID int64
		want       *Customer
		wantErr    error
	}{
		{
			name:       "when id is found return customer",
			customerID: insertedCustomerID,
			want: &Customer{
				ID:      insertedCustomerID,
				Address: "address",
				City:    "city",
				State:   "state",
			},
		},
		{
			name:       "when id is not found return error",
			customerID: -1,
			wantErr:    sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &CustomerRepository{
				DB: s.DB,
			}
			got, err := r.GetCustomer(tt.customerID)
			require.Equal(t, tt.wantErr, err)

			if tt.want != nil {
				require.Equal(t, tt.want.ID, got.ID)
				require.Equal(t, tt.want.Address, got.Address)
				require.Equal(t, tt.want.City, got.City)
				require.Equal(t, tt.want.State, got.State)
			}
		})
	}
}
