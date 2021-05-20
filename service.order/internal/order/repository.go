package order

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/postgres"
)

// Repository db repository
type Repository interface {
	// OrderPaid needs to change order status, and probably notify other services, such as product service to
	// decrease the product sku.
	OrderPaid(ctx context.Context, orderID string) error
}

type repository struct {
	pdb postgres.Provider
}

// NewRepository returns a new repostory
func NewRepository(pdb postgres.Provider) Repository {
	return &repository{pdb}
}

func (r repository) OrderPaid(ctx context.Context, orderID string) error {
	// do necessary db update to service order tables, such as updating order status.
	return nil
}
