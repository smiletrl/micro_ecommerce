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

	// echo instance
	e := echo.New()
	echoGroup := e.Group("/api/v1")

	tracingProvider := tracing.NewProvider()
	closer, err := tracingProvider.SetupTracer(constants.TracingCart, cfg)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// middleware
	e.Use(accesslog.Middleware(logger))
	e.Use(tracingProvider.Middleware(logger))
	e.Use(errors.Middleware(logger))

	// initialize service
	healthcheck.RegisterHandlers(e.Group(""))

	// redis
	redis, err := redis.NewProvider(cfg, tracingProvider)
	if err != nil {
		panic(err)
	}

	jwtProvider := jwt.NewProvider(cfg.JwtSecret)

	// product grpc client.
	pclient, err := productClient.NewClient(cfg.InternalServer.Product, tracingProvider, logger)
	if err != nil {
		panic(err)
	}

	// cart
	cartRepo := cart.NewRepository(redis, tracingProvider)
	productProxy := product{pclient}
	cartService := cart.NewService(cartRepo, productProxy, logger)
	cart.RegisterHandlers(echoGroup, cartService, jwtProvider)

	// start server
	panic(e.Start(constants.RestPort))
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
