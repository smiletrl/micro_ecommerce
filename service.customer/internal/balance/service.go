package balance

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Service balance
type Service interface {
	Add(c echo.Context, customerID int64, balance int) error
}

type service struct {
	repo   Repository
	logger *zap.SugaredLogger
}

// NewRepository returns a new repostory
func NewService(repo Repository, logger *zap.SugaredLogger) Service {
	return service{repo, logger}
}

func (s service) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
