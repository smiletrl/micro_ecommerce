package balance

import (
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/context"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"github.com/labstack/echo"
	"github.com/medivhzhan/weapp"
	"github.com/pkg/errors"
	"unicode/utf8"
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

func (s Service) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
