package main

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/balance"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/customer"
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
	closer, err := tracingProvider.SetupTracer(constants.TracingCustomer, cfg)
	if err != nil {
		logger.Fatal("tracing", err)
	}
	defer closer.Close()

	// middleware
	e.Use(accesslog.Middleware(logger))
	e.Use(tracingProvider.Middleware(logger))
	e.Use(errors.Middleware(logger))

	// initialize health
	healthcheck.RegisterHandlers(e.Group(""))

	// postgres connection
	pdb, err := postgresql.NewProvider(cfg, tracingProvider)
	if err != nil {
		logger.Fatal("postgres", err)
	}
	defer pdb.Close()

	// balance
	balance.RegisterHandlers(echoGroup)

	//err = balance.Consume(config.RocketMQ)
	//if err != nil {
	//	panic(err)
	//}

	// customer
	customerRepo := customer.NewRepository(pdb)
	customerService := customer.NewService(customerRepo, logger)
	customer.RegisterHandlers(echoGroup, customerService)

	// Start server
	e.Logger.Fatal(e.Start(constants.RestPort))
}
