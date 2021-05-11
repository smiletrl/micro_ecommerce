package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/mongodb"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.product/internal/product"
	rpcserver "github.com/smiletrl/micro_ecommerce/service.product/internal/rpc/server"
	"go.uber.org/zap"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()
	echoGroup := e.Group("/api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// init config
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	config, err := config.Load(stage)
	if err != nil {
		panic(err)
	}

	// init tracing
	tracingProvider := tracing.NewProvider()
	closer, err := tracingProvider.SetupTracer("product", config)
	if err != nil {
		sugar.Fatal(err)
	}
	defer closer.Close()

	// initialize service
	healthcheck.RegisterHandlers(e.Group(""))

	db, err := mongodb.DB(config.MongoDB)
	if err != nil {
		panic(err)
	}

	// product
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo, sugar)
	product.RegisterHandlers(echoGroup, productService)

	//err = product.Consume(config.RocketMQ)
	//if err != nil {
	//	panic(err)
	//}

	// start grpc server
	go func() {
		err = rpcserver.Register(sugar)
		if err != nil {
			panic(err)
		}
	}()

	// Start rest server
	e.Logger.Fatal(e.Start(constants.RestPort))
	//e.Logger.Fatal(e.Start(":1324"))
}
