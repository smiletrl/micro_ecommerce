package customer

import (
	"database/sql"
	"time"
)

type customer struct {
	ID        int64          `db:"id" json:"id"`
	Email     string         `db:"email" json:"email"`
	FirstName sql.NullString `db:"first_name" json:"first_name"`
	LastName  sql.NullString `db:"last_name" json:"last_name"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}
