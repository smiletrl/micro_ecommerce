package payment

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
)

// Repository db repository
type Repository interface {
	// Get payment success flag
	GetProcessedFlag(ctx context.Context, orderID string) (bool, error)

	// Update payment success flag
	SetProcessedFlag(ctx context.Context, orderID string) error

	GetPaymentMethod(ctx context.Context, orderID string) (paymentMethod, error)
}

type repository struct {
	pdb postgresql.Provider
}

// NewRepository returns a new repostory
func NewRepository(pdb postgresql.Provider) Repository {
	return &repository{pdb}
}

func (r repository) GetProcessedFlag(c context.Context, orderID string) (bool, error) {
	return true, nil
}

func (r repository) SetProcessedFlag(c context.Context, orderID string) error {
	return nil
}

type paymentMethod struct {
	CustomerID int64
	Amount     int
	Method     string
}

func (r repository) GetPaymentMethod(ctx context.Context, orderID string) (paymentMethod, error) {
	return paymentMethod{
		CustomerID: int64(12),
		Amount:     1200,
		Method:     "balance",
	}, nil
}
