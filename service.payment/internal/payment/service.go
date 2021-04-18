package payment

import (
	"github.com/labstack/echo/v4"
)

// Service is service
type Service interface {
	Create(c echo.Context, customerID, orderID int64, amount int, payType string) (err error)
}

type service struct {
	repo Repository
}

// NewService is to create new service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Create(c echo.Context, customerID, orderID int64, amount int, payType string) (err error) {
	return s.repo.Create(c, customerID, orderID, amount, payType)
}
