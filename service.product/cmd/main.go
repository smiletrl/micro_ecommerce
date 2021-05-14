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
	closer, err := tracingProvider.SetupTracer(constants.TracingProduct, cfg)
	if err != nil {
		logger.Fatal("tracing", err)
	}
	defer closer.Close()

	// middleware
	e.Use(accesslog.Middleware(logger))
	e.Use(tracingProvider.Middleware(logger))
	e.Use(errors.Middleware(logger))

	// initialize service
	healthcheck.RegisterHandlers(e.Group(""))

	db, err := mongodb.DB(cfg.MongoDB)
	if err != nil {
		panic(err)
	}

	// product
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo, logger)
	product.RegisterHandlers(echoGroup, productService)

	//err = product.Consume(config.RocketMQ)
	//if err != nil {
	//	panic(err)
	//}

	// start grpc server
	go func() {
		err = rpcserver.Register(logger)
		if err != nil {
			panic(err)
		}
	}()

	// Start rest server
	e.Logger.Fatal(e.Start(constants.RestPort))
	//e.Logger.Fatal(e.Start(":1324"))
}
