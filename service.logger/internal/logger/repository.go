package logger

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/kafka"
)

// Repository db repository
type Repository interface {
	Create(c echo.Context, customerID int64, lType, productID, category string) error
}

type repository struct {
	kafka kafka.Provider
}

// NewRepository returns a new repostory
func NewRepository(kafka kafka.Provider) Repository {
	return &repository{kafka}
}

func (r repository) Create(c echo.Context, customerID int64, lType, productID, category string) error {
	ctx := c.Request().Context()
	partition := 0
	messages := fmt.Sprintf("customerID:%d||type:%s||category:%s||productID:%s", customerID, lType, category, productID)
	err := r.kafka.Produce(ctx, constants.KafkaTopic("logger"), partition, messages)
	if err != nil {
		return errors.Wrapf(errorsd.New("error producing message"), "error producing message: %s", err.Error())
	}
	return nil
}
