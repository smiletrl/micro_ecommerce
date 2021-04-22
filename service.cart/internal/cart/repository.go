package cart

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"strconv"
)

// Repository db repository
type Repository interface {
	Get(c echo.Context, customerID int64) (items map[string]string, err error)

	Create(c echo.Context, customerID int64, skuID string, quantity int) error

	Update(c echo.Context, customerID int64, skuID string, quantity int) error

	Delete(c echo.Context, customerID int64, skuID string) error
}

type repository struct {
	rdb *redis.Client
}

// NewRepository returns a new repostory
func NewRepository(rdb *redis.Client) Repository {
	return &repository{rdb}
}

func (r repository) Get(c echo.Context, customerID int64) (items map[string]string, err error) {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	result, err := r.rdb.HGetAll(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return items, nil
		}
		return items, err
	}
	return result, nil
}

func (r repository) Create(c echo.Context, customerID int64, skuID string, quantity int) error {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	// if this sku doesn't exist, create a new hash
	if isExisting := r.rdb.HExists(context.Background(), key, skuID).Val(); !isExisting {
		_, err := r.rdb.HSet(context.Background(), key, skuID, quantity).Result()
		if err != nil {
			return err
		}
	} else {
		// increase the sku quantity in cart
		currentQuantity, err := r.rdb.HGet(context.Background(), key, skuID).Int()
		if err != nil {
			return err
		}
		_, err = r.rdb.HSet(context.Background(), skuID, currentQuantity+quantity).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r repository) Update(c echo.Context, customerID int64, skuID string, quantity int) error {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	_, err := r.rdb.HSet(context.Background(), key, skuID, quantity).Result()
	return err
}

func (r repository) Delete(c echo.Context, customerID int64, skuID string) error {
	key := fmt.Sprintf("cart:%s", strconv.FormatInt(customerID, 10))
	_, err := r.rdb.HDel(context.Background(), key, skuID).Result()
	return err
}
