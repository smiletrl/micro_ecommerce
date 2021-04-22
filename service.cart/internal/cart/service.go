package cart

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"strconv"
)

// Service is cart service
type Service interface {
	Get(c echo.Context) (cart []cartItem, err error)

	// create new cart
	Create(c echo.Context, skuID string, quantity int) error

	// update cart
	Update(c echo.Context, skuID string, quantity int) error

	// delete cart
	Delete(c echo.Context, skuID string) error
}

type service struct {
	repo         Repository
	productProxy ProductProxy
}

// NewService is to create new service
func NewService(repo Repository, product ProductProxy) Service {
	return service{repo, product}
}

func (s service) Get(c echo.Context) (cart []cartItem, err error) {
	// depending on the performance, maybe add the result to redis cache.
	customerID := c.Get(constants.AuthCustomerID).(int64)

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
		quantity, err := strconv.Atoi(items[property.SkuID])
		if err != nil {
			return cart, err
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

	return cart, err
}

func (s service) Create(c echo.Context, skuID string, quantity int) error {

	customerID := c.Get(constants.AuthCustomerID).(int64)

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

func (s service) Update(c echo.Context, skuID string, quantity int) error {

	customerID := c.Get(constants.AuthCustomerID).(int64)

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

func (s service) Delete(c echo.Context, skuID string) error {

	customerID := c.Get(constants.AuthCustomerID).(int64)

	return s.repo.Delete(c, customerID, skuID)
}
