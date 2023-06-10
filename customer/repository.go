package customer

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type IRepository interface {
	GetCustomer(customerID int64) (*Customer, error)
}

type Customer struct {
	ID        int64     `db:"id"`
	Address   string    `db:"address"`
	City      string    `db:"city"`
	State     string    `db:"state"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CustomerRepository struct {
	DB *sqlx.DB
}

func (r *CustomerRepository) GetCustomer(customerID int64) (*Customer, error) {
	var customer Customer
	err := r.DB.Get(&customer, "SELECT * FROM customer WHERE id = $1", customerID)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}
