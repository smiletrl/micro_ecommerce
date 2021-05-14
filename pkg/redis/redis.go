package redis

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"strings"
	"time"
)

type Provider interface {
	HExistsVal(ctx context.Context, key, field string) bool
	HGetInt(ctx context.Context, key, field string) (int, error)
	HGetAllResult(ctx context.Context, key string) (map[string]string, error)
	HSetResult(ctx context.Context, key string, values ...interface{}) (int64, error)
	HDelResult(ctx context.Context, key string, fields ...string) (int64, error)
}

type provider struct {
	db      *redis.Client
	tracing tracing.Provider
}

func NewProvider(cfg config.Config, tracing tracing.Provider) Provider {
	db := DB(cfg, 0)
	return provider{db, tracing}
}

func (p provider) HExistsVal(ctx context.Context, key, field string) bool {
	span, ctx := p.tracing.StartSpan(ctx, "Redis HExistsVal:"+key)
	defer p.tracing.FinishSpan(span)

	return p.db.HExists(ctx, key, field).Val()
}

func (p provider) HGetInt(ctx context.Context, key, field string) (int, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Redis HGetInt:"+key)
	defer p.tracing.FinishSpan(span)

	return p.db.HGet(ctx, key, field).Int()
}

func (p provider) HGetAllResult(ctx context.Context, key string) (map[string]string, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Redis HGetAllResult:"+key)
	defer p.tracing.FinishSpan(span)

	return p.db.HGetAll(ctx, key).Result()
}

func (p provider) HSetResult(ctx context.Context, key string, values ...interface{}) (int64, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Redis HSetResult:"+key)
	defer p.tracing.FinishSpan(span)

	return p.db.HSet(ctx, key, values...).Result()
}

func (p provider) HDelResult(ctx context.Context, key string, fields ...string) (int64, error) {
	span, ctx := p.tracing.StartSpan(ctx, "Redis HDelResult:"+key)
	defer p.tracing.FinishSpan(span)

	return p.db.HDel(ctx, key, fields...).Result()
}

// DB creates a new redis client
func DB(cfg config.Config, position int) *redis.Client {
	// redis service
	redisOptions := &redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       position,
	}
	stage := cfg.Stage
	if stage == constants.StageProd {
		redisOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	rdb := redis.NewClient(redisOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	return rdb
}

func NewMockProvider(cfg config.Config, tracing tracing.Provider) Provider {
	// mock use the second redis db instance.
	db := DB(cfg, 1)
	if strings.Contains(cfg.Stage, constants.StageGithub) {
		// in github test, use the default one
		db = DB(cfg, 0)
	}
	return provider{db, tracing}
}
