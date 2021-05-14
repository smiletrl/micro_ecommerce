package customer

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
)

// Service is cutomer service
type Service interface {
	Get(c context.Context, id int64) (cus customer, err error)

	// create new customer
	Create(c context.Context, email, firstName, lastName string) (id int64, err error)

	// update customer
	Update(c context.Context, id int64, email, firstName, lastName string) error

	// delete customer
	Delete(c context.Context, id int64) error
}

type service struct {
	repo   Repository
	logger logger.Provider
}

// NewService is to create new service
func NewService(repo Repository, logger logger.Provider) Service {
	return &service{repo, logger}
}

func (s *service) Get(c context.Context, id int64) (cus customer, err error) {
	return s.repo.Get(c, id)
}

func (s *service) Create(c context.Context, email, firstName, lastName string) (id int64, err error) {
	return s.repo.Create(c, email, firstName, lastName)
}

func (s *service) Update(c context.Context, id int64, email, firstName, lastName string) error {
	return s.repo.Update(c, id, email, firstName, lastName)
}

func (s *service) Delete(c context.Context, id int64) error {
	return s.repo.Delete(c, id)
}
