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

type provider struct {
	config  config.Config
	tracing tracing.Provider
	logger  logger.Provider
	pdb     postgresql.Provider
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
	closer, err := tracing.SetupTracer(constants.TracingCustomer, cfg)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// init postgres
	pdb, err := postgresql.NewProvider(cfg, tracing)
	if err != nil {
		panic(err)
	}
	defer pdb.Close()

	p := provider{
		config:  cfg,
		tracing: tracing,
		logger:  logger,
		pdb:     pdb,
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

	group := e.Group("/api/v1")

	// balance
	balance.RegisterHandlers(group)

	// customer
	customerRepo := customer.NewRepository(p.pdb)
	customerService := customer.NewService(customerRepo, p.logger)
	customer.RegisterHandlers(group, customerService)

	// Start server
	panic(e.Start(constants.RestPort))
}
