package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// @todo, add trace
func DB(cfg config.MongodbConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	defer func() {
		//if err = client.Disconnect(ctx); err != nil {
		//	panic(err)
		//}
	}()
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.Name), nil
}
