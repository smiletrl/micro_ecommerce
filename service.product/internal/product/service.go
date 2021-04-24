package product

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Service is cutomer service
type Service interface {
	Get(c echo.Context, id int64) (pro product, err error)

	// create new product
	Create(c echo.Context, prod product) (id string, err error)

	// update product
	Update(c echo.Context, id int64, title string, amount, stock int) error

	// delete product
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

func (s *service) Get(c echo.Context, id int64) (pro product, err error) {
	return s.repo.Get(c, id)
}

func (s *service) Create(c echo.Context, prod product) (id string, err error) {
	return s.repo.Create(c, prod)
}

func (s *service) Update(c echo.Context, id int64, title string, amount, stock int) error {
	return s.repo.Update(c, id, title, amount, stock)
}

func (s *service) Delete(c echo.Context, id int64) error {
	return s.repo.Delete(c, id)
}
