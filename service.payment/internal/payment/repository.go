package payment

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/postgres"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
)

// Repository db repository
type Repository interface {
	// Get payment success flag.
	GetProcessedFlag(ctx context.Context, orderID string) (bool, error)

	// Update payment success flag
	SetProcessedFlag(ctx context.Context, orderID string) error

	GetPaymentMethod(ctx context.Context, orderID string) (paymentMethod, error)
}

type repository struct {
	pdb postgres.Provider
	rdb redis.Provider
}

// NewRepository returns a new repostory
func NewRepository(pdb postgres.Provider, rdb redis.Provider) Repository {
	return &repository{pdb, rdb}
}

func (r repository) GetProcessedFlag(c context.Context, orderID string) (bool, error) {
	// redis
	return true, nil
}

func (r repository) SetProcessedFlag(c context.Context, orderID string) error {
	// redis
	return nil
}

type paymentMethod struct {
	CustomerID int64
	Amount     int
	Method     string
}

func (r repository) GetPaymentMethod(ctx context.Context, orderID string) (paymentMethod, error) {
	// postgres
	return paymentMethod{
		CustomerID: int64(12),
		Amount:     1200,
		Method:     "balance",
	}, nil
}
