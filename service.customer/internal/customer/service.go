package customer

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Service is cutomer service
type Service interface {
	Get(c echo.Context, id int64) (cus customer, err error)

	// create new customer
	Create(c echo.Context, email, firstName, lastName string) (id int64, err error)

	// update customer
	Update(c echo.Context, id int64, email, firstName, lastName string) error

	// delete customer
	Delete(c echo.Context, id int64) error
}

type service struct {
	repo   Repository
	logger *zap.SugaredLogger
}

// NewService is to create new service
func NewService(repo Repository, logger *zap.SugaredLogger) Service {
	return &service{repo, logger}
}

func (s *service) Get(c echo.Context, id int64) (cus customer, err error) {
	return s.repo.Get(c, id)
}

func (s *service) Create(c echo.Context, email, firstName, lastName string) (id int64, err error) {
	return s.repo.Create(c, email, firstName, lastName)
}

func (s *service) Update(c echo.Context, id int64, email, firstName, lastName string) error {
	return s.repo.Update(c, id, email, firstName, lastName)
}

func (s *service) Delete(c echo.Context, id int64) error {
	return s.repo.Delete(c, id)
}
