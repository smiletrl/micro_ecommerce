package cart

import (
	"github.com/labstack/echo/v4"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
)

// Service is cart service
type Service interface {
	Get(c echo.Context, id int64) (items []cartItem, err error)

	// create new cart
	Create(c echo.Context, customerID, skuID int64, quantity int) (item cartItem, err error)

	// update cart
	Update(c echo.Context, id int64, email, firstName, lastName string) error

	// delete cart
	Delete(c echo.Context, id int64) error
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
	return s.repo.Get(c, id)
}

func (s *service) Create(c echo.Context, customerID, skuID int64, quantity int) (item cartItem, err error) {
	// get product title, stock, amount/price
	sku, err := s.productProxy.GetSKU(c, skuID)
	if err != nil {
		return item, err
	}
	if sku.Stock < quantity {
		return item, errorsd.New("Product stock is not enough")
	}
	item = cartItem{
		ID:           int64(12),
		CustomerID:   int64(15),
		SKUID:        skuID,
		Quantity:     quantity,
		ProductTitle: sku.Title,
		Attributes:   sku.Attributes,
		IsValid:      true,
	}
	return item, nil

	//return s.repo.Create(c, customer_id, product_id, product.Title, quantity)
}

func (s *service) Update(c echo.Context, id int64, email, firstName, lastName string) error {
	return s.repo.Update(c, id, email)
}

func (s *service) Delete(c echo.Context, id int64) error {
	return s.repo.Delete(c, id)
}
