package main

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/mongodb"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.product/internal/product"
	rpcserver "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/server"
	"os"
)

type provider struct {
	config  config.Config
	tracing tracing.Provider
	logger  logger.Provider
	mongodb mongodb.Provider
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

	// mongodb connection
	mdb, err := mongodb.NewProvider(cfg.MongoDB, tracing)
	if err != nil {
		panic(err)
	}
	defer mdb.Close()

	p := provider{
		config:  cfg,
		tracing: tracing,
		logger:  logger,
		mongodb: mdb,
	}
	buildRegisters(p)

	//err = product.Consume(config.RocketMQ)
	//if err != nil {
	//	panic(err)
	//}

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

	// product
	productRepo := product.NewRepository(p.mongodb)
	productService := product.NewService(productRepo, p.logger)
	product.RegisterHandlers(group, productService)

	// start grpc server
	go func() {
		panic(rpcserver.Register(p.logger))
	}()

	// Start rest server
	panic(e.Start(constants.RestPort))
}
