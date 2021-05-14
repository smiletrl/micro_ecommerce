package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Provider interface {
	FindOne(collection string, ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult

	InsertOne(collection string, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	UpdateOne(collection string, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	DeleteOne(collection string, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	Close()
}

type provider struct {
	mdb     *mongo.Database
	mclient *mongo.Client
	tracing tracing.Provider
}

func NewProvider(cfg config.MongodbConfig, tracing tracing.Provider) (Provider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}

	// ping test
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	mdb := client.Database(cfg.Name)
	return provider{mdb, client, tracing}, nil
}

func (p provider) FindOne(collection string, ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	span, ctx := p.tracing.StartSpan(ctx, "Mongodb FindOne: "+collection)
	defer p.tracing.FinishSpan(span)

	return p.mdb.Collection(collection).FindOne(ctx, filter, opts...)
}

func (p provider) InsertOne(collection string, ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Mongodb InsertOne: "+collection)
	defer p.tracing.FinishSpan(span)

	return p.mdb.Collection(collection).InsertOne(ctx, document, opts...)
}

func (p provider) UpdateOne(collection string, ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Mongodb UpdateOne: "+collection)
	defer p.tracing.FinishSpan(span)

	return p.mdb.Collection(collection).UpdateOne(ctx, filter, update, opts...)
}

func (p provider) DeleteOne(collection string, ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Mongodb DeleteOne: "+collection)
	defer p.tracing.FinishSpan(span)

	return p.mdb.Collection(collection).DeleteOne(ctx, filter, opts...)
}

func (p provider) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.mclient.Disconnect(ctx); err != nil {
		panic(err)
	}
}
