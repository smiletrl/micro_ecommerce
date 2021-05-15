package payment

import "time"

// payment models payment
type payment struct {
	ID               int64     `db:"id"`
	Type             string    `db:"type"`
	OrderID          int64     `db:"order_id"`
	CustomerID       int64     `db:"customer_id"`
	PrepayID         int64     `db:"prepay_id"`
	Amount           int       `db:"amount"`
	IsPaid           bool      `db:"is_paid"`
	IsStaffNotified  bool      `db:"is_staff_notified"`
	IsAuthorNotified bool      `db:"is_author_notified"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
