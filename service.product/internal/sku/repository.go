package sku

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/mongodb"
)

// Repository db repository
type Repository interface {
	Get(c context.Context, productID string) (skus []sku, err error)

	Create(c context.Context, s sku) error

	Update(c context.Context, id string, s sku) error

	Delete(c context.Context, id string) error
}

type repository struct {
	mdb mongodb.Provider
}

// NewRepository returns a new repostory
func NewRepository(mdb mongodb.Provider) Repository {
	return &repository{mdb}
}

func (r repository) Get(c context.Context, productID string) (skus []sku, err error) {
	cursor, err := r.mdb.Find("sku", c, bson.M{"productId": productID})
	if err != nil {
		return skus, errorsd.New("error finding sku from db", "error finding sku from db: %s", err.Error())
	}

	var skuMs []bson.M
	if err = cursor.All(c, &skuMs); err != nil {
		return skus, errorsd.New("error cursoring sku from db", "error cursoring sku from db: %s", err.Error())
	}

	bsonBytes, _ := bson.Marshal(skuMs)
	bson.Unmarshal(bsonBytes, &skus)
	return skus, nil
}

func (r repository) Create(c context.Context, s sku) error {
	// @todo add sku validation
	_, err := r.mdb.InsertOne("sku", c, bson.D{
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

func (r repository) Update(c context.Context, id string, s sku) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.mdb.UpdateOne(
		"sku",
		c,
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

func (r repository) Delete(c context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.mdb.DeleteOne(
		"sku",
		c,
		bson.M{"_id": objectID})
	if err != nil {
		return errors.Wrapf(errorsd.New("error deleting sku in db"), "error deleting sku in db: %s", err.Error())
	}
	return nil
}
