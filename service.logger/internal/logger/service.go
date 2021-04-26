package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"go.uber.org/zap"
)

// Service is logger service
type Service interface {
	Create(c echo.Context, req createRequest) error
}

type service struct {
	repo   Repository
	logger *zap.SugaredLogger
}

// NewService is to create new service
func NewService(repo Repository, logger *zap.SugaredLogger) Service {
	return service{repo, logger}
}

func (s service) Create(c echo.Context, req createRequest) error {
	customerID := c.Get(constants.AuthCustomerID).(int64)
	return s.repo.Create(c, customerID, req.Type, req.ProductID, req.Category)
}
