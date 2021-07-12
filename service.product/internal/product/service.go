package product

import (
	"context"

	"github.com/smiletrl/micro_ecommerce/pkg/logger"
)

// Service is cutomer service
type Service interface {
	Get(ctx context.Context, id string) (pro product, err error)

	// create new product
	Create(ctx context.Context, req createRequest) (id string, err error)

	// update product
	Update(ctx context.Context, id string, req updateRequest) error

	// delete product
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo   Repository
	logger logger.Provider
}

// NewService is to create new service
func NewService(repo Repository, logger logger.Provider) Service {
	return &service{repo, logger}
}

func (s *service) Get(ctx context.Context, id string) (pro product, err error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, req createRequest) (id string, err error) {
	return s.repo.Create(ctx, req)
}

func (s *service) Update(ctx context.Context, id string, req updateRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
