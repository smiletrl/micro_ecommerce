package cart

import (
	"github.com/labstack/echo/v4"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
)

// Service is cart service
type Service interface {
	Get(c echo.Context, customerID int64) (items []cartItem, err error)

	// create new cart
	Create(c echo.Context, customerID, skuID int64, quantity int) (item cartItem, err error)

	// update cart
	Update(c echo.Context, customerID string, skuID int64, quantity int) error

	// delete cart
	Delete(c echo.Context, customerID string, skuID int64) error
}

type service struct {
	repo         Repository
	productProxy ProductProxy
}

// NewService is to create new service
func NewService(repo Repository, product ProductProxy) Service {
	return &service{repo, product}
}

func (s *service) Get(c echo.Context, id int64) (items []cartItem, err error) {
	// @todo get each cart item detail from product service
	// depending on the performance, maybe add the result to redis cache too.
	return s.repo.Get(c, id)
}

func (s *service) Create(c echo.Context, customerID, skuID int64, quantity int) (item cartItem, err error) {
	// get product sku stock & verify this cart item can be created.
	sku, err := s.productProxy.GetSKU(c, skuID)
	if err != nil {
		return item, err
	}
	if sku.Stock < quantity {
		return item, errorsd.New("Product stock is not enough")
	}
	item = cartItem{
		CustomerID: int64(15),
		SkuID:      "123344qwqw",
	}
	return item, nil

	//return s.repo.Create(c, customer_id, product_id, product.Title, quantity)
}

func (s *service) Update(c echo.Context, customerID string, skuID int64, quantity int) error {
	return nil
	//return s.repo.Update(c, id, email)
}

func (s *service) Delete(c echo.Context, customerID string, skuID int64) error {
	return nil
	//return s.repo.Delete(c, id)
}
