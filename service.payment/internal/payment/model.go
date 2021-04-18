package payment

import "time"

// payment models payment
type payment struct {
	ID               int64     `db:"id"`
	Type             string    `db:"type"`
	OrderID          int64     `db:"order_id"`
	PrepayID         int64     `db:"prepay_id"`
	Amount           int       `db:"amount"`
	IsPaid           bool      `is_paid`
	IsStaffNotified  bool      `is_staff_notified`
	IsAuthorNotified bool      `is_author_notified`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
