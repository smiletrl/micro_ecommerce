package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	_ "github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.cart/internal/cart"
	productClient "github.com/smiletrl/micro_ecommerce/service.product/external"

	"os"
)

type provider struct {
	config  config.Config
	tracing tracing.Provider
	logger  logger.Provider
	jwt     jwt.Provider
	rdb     redis.Provider
}

func main() {
	// stage
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}

	// init config
	cfg, err := config.Load(stage)
	if err != nil {
		panic(err)
	}

	// init logger
	logger := logger.NewProvider(cfg.Logger)
	defer logger.Close()

	// init tracing
	tracing, err := tracing.NewProvider(constants.TracingCart, cfg)
	if err != nil {
		panic(err)
	}
	defer tracing.Close()

	// redis
	redis, err := redis.NewProvider(cfg, tracing)
	if err != nil {
		panic(err)
	}

	jwtProvider := jwt.NewProvider(cfg.JwtSecret)

	p := provider{
		config:  cfg,
		tracing: tracing,
		logger:  logger,
		jwt:     jwtProvider,
		rdb:     redis,
	}
	buildRegisters(p)
}

// product proxy
type product struct {
	client productClient.Client
}

func (p product) GetSkuStock(c context.Context, skuID string) (int, error) {
	return p.client.GetSkuStock(c, skuID)
}

func (p product) GetSkuProperties(c context.Context, skuIDs []string) ([]entity.SkuProperty, error) {
	return p.client.GetSkuProperties(c, skuIDs)
}

func buildRegisters(p provider) {
	// echo instance
	e := echo.New()

	// middleware
	e.Use(accesslog.Middleware(p.logger))
	e.Use(p.tracing.Middleware(p.logger))
	e.Use(errors.Middleware(p.logger))

	// initialize health
	healthcheck.RegisterHandlers(e.Group(""))

	group := e.Group("/api/v1")

	// cart
	cartRepo := cart.NewRepository(p.rdb, p.tracing)

	// product grpc client.
	pclient, err := productClient.NewClient(p.config.InternalServer.Product, p.tracing, p.logger)
	if err != nil {
		panic(err)
	}

	productProxy := product{pclient}
	cartService := cart.NewService(cartRepo, productProxy, p.logger)
	cart.RegisterHandlers(group, cartService, p.jwt)

	// start server
	panic(e.Start(constants.RestPort))
}
