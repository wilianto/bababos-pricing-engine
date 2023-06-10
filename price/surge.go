package price

import "github.com/wilianto/bababos-pricing-engine/customer"

type PricingSurgeStrategy struct {
	PriceRepository    IRepository
	CustomerRepository customer.IRepository
}

func (s *PricingSurgeStrategy) GetPrice(req PriceRequest) (PriceResponse, error) {
	// TODO: implement pricing with surge strategy
	// This one will be presented in the presentation, not gonna implement it here
	return PriceResponse{}, nil
}
