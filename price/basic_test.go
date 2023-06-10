package price

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wilianto/bababos-pricing-engine/customer"
)

type MockPriceRepository struct {
}

func (r *MockPriceRepository) FindLowestPrice(skuID string, state string) (*RecommendedPrice, error) {
	return &RecommendedPrice{
		Price:      1000,
		SupplierID: 1,
	}, nil
}

type MockCustomerRepository struct {
}

func (r *MockCustomerRepository) GetCustomer(customerID int64) (*customer.Customer, error) {
	return &customer.Customer{
		ID:      1,
		Address: "Jalan Raya",
		City:    "Jakarta",
		State:   "Jakarta",
	}, nil
}

type PricingBasicStrategyTestSuite struct {
	suite.Suite
}

func TestPricingBasicStrategyTestSuite(t *testing.T) {
	suite.Run(t, new(PricingBasicStrategyTestSuite))
}

func (s *PricingBasicStrategyTestSuite) TestGetPrice() {
	tests := []struct {
		name           string
		qty            int64
		suggestedPrice float64
		wantErr        error
	}{
		{
			name:           "when qty is more than 100 then mark up 10%",
			qty:            150,
			suggestedPrice: 1100,
		},
		{
			name:           "when qty is 80 then mark up 15%",
			qty:            80,
			suggestedPrice: 1150,
		},
		{
			name:           "when qty is 30 then mark up 20%",
			qty:            30,
			suggestedPrice: 1200,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := &PricingBasicStrategy{
				PriceRepository:    &MockPriceRepository{},
				CustomerRepository: &MockCustomerRepository{},
			}
			got, err := r.GetPrice(PriceRequest{
				SkuID:      "123456",
				CustomerID: 1,
				Qty:        tt.qty,
				Unit:       "kg",
			})
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.suggestedPrice, got.SuggestedPrice)
		})
	}
}
