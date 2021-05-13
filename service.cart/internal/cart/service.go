package cart

import (
	"context"
	_ "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"strconv"
)

// Service is cart service
type Service interface {
	Get(c context.Context) (cart []cartItem, err error)

	// create new cart
	Create(c context.Context, skuID string, quantity int) error

	// update cart
	Update(c context.Context, skuID string, quantity int) error

	// delete cart
	Delete(c context.Context, skuID string) error
}

type service struct {
	repo         Repository
	productProxy ProductProxy
	log          logger.Provider
}

// NewService is to create new service
func NewService(repo Repository, product ProductProxy, log logger.Provider) Service {
	return service{repo, product, log}
}

func (s service) Get(c context.Context) (cart []cartItem, err error) {
	customerID := c.Value(constants.AuthCustomerID).(int64)

	// items is a map, key is skuID, value is this sku's quantity in cart
	items, err := s.repo.Get(c, customerID)
	if err != nil {
		return cart, err
	}

	// get all sku ids and retrieve the sku properties
	skuIDs := make([]string, len(items))
	i := 0
	for skuID := range items {
		skuIDs[i] = skuID
		i++
	}

	// get sku properties from product service
	properties, err := s.productProxy.GetSkuProperties(c, skuIDs)
	if err != nil {
		return cart, err
	}

	cart = make([]cartItem, len(items))
	j := 0

	for _, property := range properties {
		// double check before using this sku property
		rawQuantity, ok := items[property.SkuID]
		if !ok {
			continue
		}
		quantity, err := strconv.Atoi(rawQuantity)
		if err != nil {
			return cart, errors.Wrapf(errorsd.New("error converting sku quantity to int"), "error converting sku quantity: %s", err.Error())
		}
		valid := true
		// if sku stock is less than sku quantity in cart, then this cart item
		// will be no longer valid.
		if property.Stock < quantity {
			valid = false
		}
		cart[j] = cartItem{
			Quantity:    quantity,
			Valid:       valid,
			SkuProperty: property,
		}
		j++
	}

	// @todo, add the deleted sku id item to response.
	// It's possible that one sku from cart has been deleted in system, then
	// cart should probably tell user this cart item doesn't exist any more.

	return cart, err
}

func (s service) Create(c context.Context, skuID string, quantity int) error {
	customerID := c.Value(constants.AuthCustomerID).(int64)

	// get product sku stock
	stock, err := s.productProxy.GetSkuStock(c, skuID)
	if err != nil {
		return err
	}
	if stock < quantity {
		return errorsd.New("product stock is not enough")
	}
	return s.repo.Create(c, customerID, skuID, quantity)
}

func (s service) Update(c context.Context, skuID string, quantity int) error {
	customerID := c.Value(constants.AuthCustomerID).(int64)

	if quantity < 1 {
		return errorsd.New("quantity can not be under 1")
	}

	// get product sku stock
	stock, err := s.productProxy.GetSkuStock(c, skuID)
	if err != nil {
		return err
	}
	if stock < quantity {
		return errorsd.New("product stock is not enough")
	}

	return s.repo.Update(c, customerID, skuID, quantity)
}

func (s service) Delete(c context.Context, skuID string) error {
	customerID := c.Value(constants.AuthCustomerID).(int64)

	return s.repo.Delete(c, customerID, skuID)
}
