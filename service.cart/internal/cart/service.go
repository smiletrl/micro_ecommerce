package cart

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
)

// Service is cart service
type Service interface {
	Get(c echo.Context, id int64) (items []cartItem, err error)

	// create new customer
	Create(c echo.Context, customer_id, product_id int64, quantity int) (id int64, err error)

	// update customer
	Update(c echo.Context, id int64, email, firstName, lastName string) error

	// delete customer
	Delete(c echo.Context, id int64) error
}

type ProductProxy interface {
	GetDetail(c echo.Context, id int64) (entity.Product, error)
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

func (s *service) Create(c echo.Context, customer_id, product_id int64, quantity int) (id int64, err error) {
	// create the cart items
	// get product title, stock, amount/price
	product, _ := s.productProxy.GetDetail(c, product_id)
	if product.Stock < quantity {
		return id, errorsd.New("Product stock is not enough")
	}
	return s.repo.Create(c, customer_id, product_id, product.Title, quantity)
}

func (s *service) Update(c echo.Context, id int64, email, firstName, lastName string) error {
	return s.repo.Update(c, id, email)
}

func (s *service) Delete(c echo.Context, id int64) error {
	return s.repo.Delete(c, id)
}
