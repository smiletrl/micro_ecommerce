package cache

import (
	"context"
	"crypto/tls"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/go-redis/redis/v8"
	"log"
)

// NewRedis new redis client
func NewRedis(cfg *config.Config) *redis.Client {
	// redis service
	redisOptions := &redis.Options{
		Addr:     cfg.Redis.Endpoint + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0, // use default DB
	}
	stage := cfg.Stage
	if stage != constants.StageLocal && stage != constants.StageGithub {
		redisOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	rdb := redis.NewClient(redisOptions)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("redis ping error: %s", err.Error())
		panic(err)
	}
	return rdb
}
