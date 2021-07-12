package cart

import (
	"context"
	"fmt"
	"strconv"

	goredis "github.com/go-redis/redis/v8"

	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

// Repository db repository
type Repository interface {
	Get(c context.Context, customerID int64) (items map[string]string, err error)

	Create(c context.Context, customerID int64, skuID string, quantity int) error

	Update(c context.Context, customerID int64, skuID string, quantity int) error

	Delete(c context.Context, customerID int64, skuID ...string) error
}

type repository struct {
	rdb     redis.Provider
	tracing tracing.Provider
}

// NewRepository returns a new repostory
func NewRepository(rdb redis.Provider, tracing tracing.Provider) Repository {
	return &repository{rdb, tracing}
}

func (r repository) Get(c context.Context, customerID int64) (items map[string]string, err error) {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	result, err := r.rdb.HGetAllResult(c, key)
	if err != nil {
		if err == goredis.Nil {
			return items, nil
		}
		return items, err
	}
	return result, nil
}

func (r repository) Create(c context.Context, customerID int64, skuID string, quantity int) error {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	// if this sku doesn't exist, create a new hash
	if isExisting := r.rdb.HExistsVal(c, key, skuID); !isExisting {
		_, err := r.rdb.HSetResult(c, key, skuID, quantity)
		if err != nil {
			return err
		}
	} else {
		// increase the sku quantity in cart
		currentQuantity, err := r.rdb.HGetInt(c, key, skuID)
		if err != nil {
			return err
		}
		_, err = r.rdb.HSetResult(c, key, skuID, currentQuantity+quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r repository) Update(c context.Context, customerID int64, skuID string, quantity int) error {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	_, err := r.rdb.HSetResult(c, key, skuID, quantity)
	return err
}

func (r repository) Delete(c context.Context, customerID int64, skuID ...string) error {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	var err error
	if len(skuID) == 0 {
		return errorsd.New("at least one sku id is required")
	}
	_, err = r.rdb.HDelResult(c, key, skuID...)
	return err
}
