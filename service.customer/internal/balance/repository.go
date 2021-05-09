package balance

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
)

// Repository db repository
type Repository interface {

	// create new customer
	Add(c echo.Context, customerID int64, balance int) error
}

type repository struct {
	pdb postgresql.DB
}

// NewRepository returns a new repostory
func NewRepository(db postgresql.DB) Repository {
	return &repository{db}
}

func (r *repository) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
