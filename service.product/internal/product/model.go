package product

import "time"

// product is a simple model, without sku/attribute/category.
type product struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	// product's price amount
	Amount int `db:"amount"`
	// product's stock value
	Stock     int       `db:"stock"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
