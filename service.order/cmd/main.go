package main

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	_ "github.com/smiletrl/micro_ecommerce/pkg/postgresql"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.order/internal/order"
	"os"
)

type provider struct {
	config  config.Config
	tracing tracing.Provider
	logger  logger.Provider
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
	tracing := tracing.NewProvider()
	closer, err := tracing.SetupTracer(constants.TracingPayment, cfg)
	if err != nil {
		logger.Fatal("tracing", err)
	}
	defer closer.Close()

	// rocketMQ message
	err = order.Consume(cfg.RocketMQ)
	if err != nil {
		panic(err)
	}

	p := provider{
		config:  cfg,
		tracing: tracing,
		logger:  logger,
	}
	buildRegisters(p)
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

	//group := e.Group("/api/v1")

	// Start rest server
	panic(e.Start(constants.RestPort))
}
