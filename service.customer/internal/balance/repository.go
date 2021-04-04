package balance

import (
	"github.com/labstack/echo"
	"time"

	"github.com/smiletrl/micro_ecommerce/pkg/context"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// Repository db repository
type Repository interface {

	// create new customer
	Add(c echo.Context, customerID int64, balance int) error
}

type repository struct {
	db dbcontext.DB
}

// NewRepository returns a new repostory
func NewRepository(db dbcontext.DB) Repository {
	return &repository{db}
}

func (r *repository) Add(c echo.Context, customerID int64, balance int) error {
	return nil
}
