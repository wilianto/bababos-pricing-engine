package price

import "github.com/jmoiron/sqlx"

type IRepository interface {
	FindLowestPrice(skuID string, state string) (*RecommendedPrice, error)
}

type PriceRepository struct {
	DB *sqlx.DB
}

type RecommendedPrice struct {
	Price      float64 `db:"price"`
	SupplierID int64   `db:"supplier_id"`
}

func (r *PriceRepository) FindLowestPrice(skuID string, state string) (*RecommendedPrice, error) {
	var price RecommendedPrice
	query := `
		SELECT price, supplier_id 
		FROM price 
		INNER JOIN supplier ON price.supplier_id = supplier.id
		WHERE price.sku_id = $1 AND supplier.state = $2 
		ORDER BY price ASC LIMIT 1
	`
	err := r.DB.Get(&price, query, skuID, state)
	if err != nil {
		return nil, err
	}

	return &price, nil
}
