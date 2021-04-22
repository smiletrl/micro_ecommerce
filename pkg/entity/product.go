package entity

import "time"

// product is a simple model, without sku/attribute/category.
type Product struct {
	ID    int64  `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	// product's price amount
	Amount int `db:"amount" json:"amount"`
	// product's stock value
	Stock     int       `db:"stock" json:"stock"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// SKU is product sku
type SKU struct {
	Stock      int
	Amount     int
	Title      string
	Attributes string
}

type SkuProperty struct {
	SkuID      string `json:"sku_id"`
	Title      string `json:"title"`
	Price      int    `json:"price"`
	Attributes string `json:"attributes"`
	Thumbnail  string `json:"thumbnail"`
	Stock      int    `json:"stock"`
}
