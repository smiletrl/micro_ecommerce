package redis

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"time"
)

// New creates a new redis client
func New(cfg config.Config) *redis.Client {
	// redis service
	redisOptions := &redis.Options{
		Addr:     cfg.Redis.Endpoint + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0,
	}
	stage := cfg.Stage
	if stage != constants.StageLocal && stage != constants.StageGithub && stage != constants.StageK8s {
		redisOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	rdb := redis.NewClient(redisOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	return rdb
}
