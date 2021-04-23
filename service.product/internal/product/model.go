package product

import "time"

// product models product
type product struct {
	Title       string `db:"title"`
	Description string `db:"description"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
