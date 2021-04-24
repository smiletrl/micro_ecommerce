package product

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository db repository
type Repository interface {
	// get product
	Get(c echo.Context, id int64) (prod product, err error)

	// create new product
	Create(c echo.Context, prod product) (id string, err error)

	// update product
	Update(c echo.Context, id int64, title string, amount, stock int) error

	// delete product
	Delete(c echo.Context, id int64) error
}

type repository struct {
	db *mongo.Database
}

// NewRepository returns a new repostory
func NewRepository(db *mongo.Database) Repository {
	return &repository{db}
}

func (r repository) Get(c echo.Context, id int64) (pro product, err error) {
	return pro, err
}

func (r repository) Create(c echo.Context, prod product) (id string, err error) {
	// @todo add product validation
	collection := r.db.Collection("product")
	ctx := c.Request().Context()
	res, err := collection.InsertOne(ctx, bson.D{
		{"title", prod.Title},
		{"body", prod.Body},
		{"category", prod.Category},
		{"variants", prod.Variants}})
	if err != nil {
		return id, err
	}
	id = res.InsertedID.(string)
	return id, nil
}

func (r repository) Update(c echo.Context, id int64, title string, amount, stock int) error {
	return nil
}

func (r repository) Delete(c echo.Context, id int64) error {
	return nil
}
