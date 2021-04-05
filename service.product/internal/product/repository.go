package product

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
)

// Repository db repository
type Repository interface {
	// get product
	Get(c echo.Context, id int64) (pro product, err error)

	// create new product
	Create(c echo.Context, title string, amount, stock int) (id int64, err error)

	// update product
	Update(c echo.Context, id int64, title string, amount, stock int) error

	// delete product
	Delete(c echo.Context, id int64) error
}

type repository struct {
	db dbcontext.DB
}

// NewRepository returns a new repostory
func NewRepository(db dbcontext.DB) Repository {
	return &repository{db}
}

func (r repository) Get(c echo.Context, id int64) (pro product, err error) {
	return pro, err
}

func (r repository) Create(c echo.Context, title string, amount, stock int) (id int64, err error) {
	return id, nil
}

func (r repository) Update(c echo.Context, id int64, title string, amount, stock int) error {
	return nil
}

func (r repository) Delete(c echo.Context, id int64) error {
	return nil
}
