package sku

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository db repository
type Repository interface {
	Get(c echo.Context, productID string) (skus []sku, err error)

	Create(c echo.Context, s sku) error

	Update(c echo.Context, id string, s sku) error

	Delete(c echo.Context, id string) error
}

type repository struct {
	mdb *mongo.Database
}

// NewRepository returns a new repostory
func NewRepository(mdb *mongo.Database) Repository {
	return &repository{mdb}
}

func (r repository) Get(c echo.Context, productID string) (skus []sku, err error) {
	collection := r.mdb.Collection("sku")

	ctx := c.Request().Context()
	cursor, err := collection.Find(ctx, bson.M{"productId": productID})
	if err != nil {
		return skus, errorsd.New("error finding sku from db", "error finding sku from db: %s", err.Error())
	}

	var skuMs []bson.M
	if err = cursor.All(ctx, &skuMs); err != nil {
		return skus, errorsd.New("error cursoring sku from db", "error cursoring sku from db: %s", err.Error())
	}

	bsonBytes, _ := bson.Marshal(skuMs)
	bson.Unmarshal(bsonBytes, &skus)
	return skus, nil
}

func (r repository) Create(c echo.Context, s sku) error {
	// @todo add sku validation
	collection := r.mdb.Collection("sku")
	ctx := c.Request().Context()
	_, err := collection.InsertOne(ctx, bson.D{
		{"productId", s.ProductID},
		{"assets", s.Assets},
		{"attrs", s.Attrs},
		{"price", s.Price},
		{"stock", s.Stock}})
	if err != nil {
		return errors.Wrapf(errorsd.New("error inserting sku in db"), "error inserting sku in db: %s", err.Error())
	}
	return nil
}

func (r repository) Update(c echo.Context, id string, s sku) error {
	collection := r.mdb.Collection("sku")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{
				{"productId", s.ProductID},
				{"assets", s.Assets},
				{"attrs", s.Attrs},
				{"price", s.Price},
				{"stock", s.Stock}}}})
	if err != nil {
		return errors.Wrapf(errorsd.New("error updating sku in db"), "error updating sku in db: %s", err.Error())
	}
	return nil
}

func (r repository) Delete(c echo.Context, id string) error {
	collection := r.mdb.Collection("sku")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	_, err = collection.DeleteOne(
		ctx,
		bson.M{"_id": objectID})
	if err != nil {
		return errors.Wrapf(errorsd.New("error deleting sku in db"), "error deleting sku in db: %s", err.Error())
	}
	return nil
}
