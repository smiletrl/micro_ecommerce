package sku

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository db repository
type Repository interface {
	Get(c echo.Context, productID string) error

	Create(c echo.Context, productID string) error

	Update(c echo.Context, skuID string) error

	Delete(c echo.Context, skuID string) error
}

type repository struct {
	db *mongo.Database
}

// NewRepository returns a new repostory
func NewRepository(db *mongo.Database) Repository {
	return &repository{db}
}

func (r repository) Get(c echo.Context, productID string) error {
	return nil
}

func (r repository) Create(c echo.Context, productID string) error {
	return nil
}

func (r repository) Update(c echo.Context, skuID string) error {
	return nil
}

func (r repository) Delete(c echo.Context, skuID string) error {
	return nil
}
