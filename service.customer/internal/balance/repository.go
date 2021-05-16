package balance

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
)

// Repository db repository
type Repository interface {
	// Increase balance
	Increase(ctx context.Context, customerID int64, amount int) error

	// Decrease balance
	Decrease(ctx context.Context, customerID int64, amount int) error
}

type repository struct {
	pdb postgresql.Provider
}

// NewRepository returns a new repostory
func NewRepository(pdb postgresql.Provider) Repository {
	return &repository{pdb}
}

func (r repository) Increase(ctx context.Context, customerID int64, amount int) error {
	return nil
}

func (r repository) Decrease(ctx context.Context, customerID int64, amount int) error {
	return nil
}
