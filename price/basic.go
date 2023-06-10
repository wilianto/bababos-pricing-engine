package price

import (
	"github.com/wilianto/bababos-pricing-engine/customer"
)

type PricingBasicStrategy struct {
	PriceRepository    IRepository
	CustomerRepository customer.IRepository
}

func (s *PricingBasicStrategy) GetPrice(req PriceRequest) (PriceResponse, error) {
	// Get customer city
	customer, err := s.CustomerRepository.GetCustomer(req.CustomerID)
	if err != nil {
		return PriceResponse{}, err
	}

	// Get lowest supplier price
	base, err := s.PriceRepository.FindLowestPrice(req.SkuID, customer.State)
	if err != nil {
		return PriceResponse{}, err
	}

	// Calculate price
	suggestedPrice := base.Price
	if req.Qty > 100 {
		suggestedPrice = base.Price * 1.1
	} else if req.Qty > 50 {
		suggestedPrice = base.Price * 1.15
	} else {
		suggestedPrice = base.Price * 1.20
	}

	return PriceResponse{
		SuggestedPrice:        suggestedPrice,
		RecommendedSupplierID: base.SupplierID,
	}, nil
}
