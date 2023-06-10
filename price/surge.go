package price

import "github.com/wilianto/bababos-pricing-engine/customer"

type SurgePricingStrategy struct {
	PriceRepository    IRepository
	CustomerRepository customer.IRepository
}

func (s *SurgePricingStrategy) GetPrice(req PriceRequest) (PriceResponse, error) {
	// TODO: implement pricing with surge strategy
	// This one will be presented in the presentation, not gonna implement it here
	return PriceResponse{}, nil
}
