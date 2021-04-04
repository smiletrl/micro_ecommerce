package balance

import (
	"github.com/labstack/echo/v4"
)

// Service balance
type Service interface {
	Add(c echo.Context, customerID int64, balance int) error
}

type service struct {
	repo Repository
}

// NewRepository returns a new repostory
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
