package order

import (
	"context"

	"github.com/smiletrl/micro_ecommerce/pkg/logger"
)

// Service is cutomer service
type Service interface {
	CreateFromCart(c context.Context, customerID int64, req createFromCartRequest) (id int64, err error)
	CreateFromProduct(c context.Context, customerID int64, req createFromProductRequest) (id int64, err error)
}

type service struct {
	repo   Repository
	logger logger.Provider
}

// NewService is to create new service
func NewService(repo Repository, logger logger.Provider) Service {
	return &service{repo, logger}
}

func (s *service) CreateFromCart(c context.Context, customerID int64, req createFromCartRequest) (id int64, err error) {
	return int64(12), err
}

func (s *service) CreateFromProduct(c context.Context, customerID int64, req createFromProductRequest) (id int64, err error) {
	return int64(12), err
}
