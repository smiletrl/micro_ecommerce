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
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	_ "github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.cart/internal/cart"
	productClient "github.com/smiletrl/micro_ecommerce/service.product/external"
	"go.uber.org/zap"
	"os"
)

func main() {
	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// stage
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	cfg, err := config.Load(stage)
	if err != nil {
		panic(err)
	}

	// echo instance
	e := echo.New()
	echoGroup := e.Group("/api/v1")

	tracingProvider := tracing.NewProvider()
	closer, err := tracingProvider.SetupTracer("cart", cfg)
	if err != nil {
		sugar.Fatal(err)
	}
	defer closer.Close()

	// middleware
	e.Use(accesslog.Middleware(sugar))
	e.Use(tracingProvider.Middleware(sugar))
	e.Use(errors.Recover(sugar))

	// initialize service
	healthcheck.RegisterHandlers(e.Group(""))

	// redis
	redisClient := redis.New(cfg)

	jwtService := jwt.NewProvider(cfg.JwtSecret)

	// product grpc client.
	pclient, err := productClient.NewClient(cfg.InternalServer.Product, tracingProvider, sugar)
	if err != nil {
		panic(err)
	}

	// cart
	cartRepo := cart.NewRepository(redisClient, tracingProvider)

	productProxy := product{pclient}
	cartService := cart.NewService(cartRepo, productProxy, sugar)
	cart.RegisterHandlers(echoGroup, cartService, jwtService)

	// start server
	sugar.Fatal(e.Start(constants.RestPort))
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
