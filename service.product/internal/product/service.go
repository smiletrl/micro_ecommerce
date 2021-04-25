package product

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Service is cutomer service
type Service interface {
	Get(c echo.Context, id string) (pro product, err error)

	// create new product
	Create(c echo.Context, req createRequest) (id string, err error)

	// update product
	Update(c echo.Context, id string, req updateRequest) error

	// delete product
	Delete(c echo.Context, id string) error
}

type service struct {
	repo   Repository
	logger *zap.SugaredLogger
}

// NewService is to create new service
func NewService(repo Repository, logger *zap.SugaredLogger) Service {
	return &service{repo, logger}
}

func (s *service) Get(c echo.Context, id string) (pro product, err error) {
	return s.repo.Get(c, id)
}

func (s *service) Create(c echo.Context, req createRequest) (id string, err error) {
	return s.repo.Create(c, req)
}

func (s *service) Update(c echo.Context, id string, req updateRequest) error {
	return s.repo.Update(c, id, req)
}

func (s *service) Delete(c echo.Context, id string) error {
	return s.repo.Delete(c, id)
}
