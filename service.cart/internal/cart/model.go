package cart

import "time"

type cartItem struct {
	CustomerID int64     `json:"customer_id"`
	SkuID      string    `json:"sku_id"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type cartItemDetail struct {
	cartItem
	Title      string `json:"title"`
	Attributes string `json:"attributes"`
	Thumbnail  string `json:"thumbnail"`
	Valid      bool   `json:"valid"`
}
