package price

type PriceRequest struct {
	CustomerID int64  `json:"customer_id"`
	SkuID      string `json:"sku_id"`
	Qty        int64  `json:"qty"`
	Unit       string `json:"unit"`
}

type PriceResponse struct {
	SuggestedPrice        float64 `json:"suggested_price"`
	RecommendedSupplierID int64   `json:"recommended_supplier_id"`
}

type PriceStrategy interface {
	GetPrice(req PriceRequest) (PriceResponse, error)
}
