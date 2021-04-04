package balance

import (
	"github.com/labstack/echo"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
)

// Service balance
type Service interface {
	Add(c echo.Context, customerID int64, balance int) error
}

type service struct {
	db dbcontext.DB
}

// NewRepository returns a new repostory
func NewService(db dbcontext.DB) Service {
	return service{db}
}

func (s service) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
