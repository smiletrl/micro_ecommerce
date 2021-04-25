package balance

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/postgre"
)

// Repository db repository
type Repository interface {

	// create new customer
	Add(c echo.Context, customerID int64, balance int) error
}

type repository struct {
	pdb postgre.DB
}

// NewRepository returns a new repostory
func NewRepository(db postgre.DB) Repository {
	return &repository{db}
}

func (r *repository) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
