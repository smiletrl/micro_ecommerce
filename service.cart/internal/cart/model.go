package cart

import "time"

type cartItem struct {
	ID           int64     `db:"id" json:"id"`
	CustomerID   int64     `db:"customer_id" json:"customer_id"`
	SKUID        int64     `db:"sku_id" json:"sku_id"`
	ProductTitle string    `db:"product_title" json:"product_title"`
	Amount       int       `db:"amount" json:"amount"`
	Quantity     int       `db:"quantity" json:"quantity"`
	Attributes   string    `db:"attributes" json:"attributes"`
	IsValid      bool      `db:"is_valid" json:"is_valid"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
