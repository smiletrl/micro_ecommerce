package payment

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
)

// Repository db repository
type Repository interface {
	Create(c echo.Context, customerID, orderID int64, amount int, payType string) error
}

type repository struct {
	db dbcontext.DB
}

// NewRepository returns a new repostory
func NewRepository(db dbcontext.DB) Repository {
	return repository{db}
}

func (r repository) Create(c echo.Context, customerID, orderID int64, amount int, payType string) error {
	return nil
}
