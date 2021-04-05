package cart

import "time"

type cartItem struct {
	ID           int64     `db:"id"`
	CustomerID   int64     `db:"customer_id"`
	ProductID    int64     `db:"product_id"`
	ProductTitle string    `db:"product_title"`
	Amount       int       `db:"amount"`
	Quantity     int       `db:"quantity"`
	IsValid      bool      `db:"is_valid"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
