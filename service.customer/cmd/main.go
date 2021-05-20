package main

import (
	"context"
	rocketmqLib "github.com/apache/rocketmq-client-go/v2"
	rocketConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/postgres"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/balance"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/customer"
	"os"
	"time"
)

type provider struct {
	config   config.Config
	tracing  tracing.Provider
	logger   logger.Provider
	consumer rocketmqLib.PushConsumer
	rocketmq rocketmq.Provider
	pdb      postgres.Provider
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

	// init postgres
	pdb, err := postgres.NewProvider(cfg, tracing)
	if err != nil {
		panic(err)
	}
	defer pdb.Close()

	// init rocketmq
	rocketmqProvider := rocketmq.NewProvider(cfg.RocketMQ, pdb, constants.ServiceCustomer)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	consumer, err := rocketmqProvider.CreatePushConsumer(ctx, constants.RocketMQGroupPayment, rocketConsumer.Clustering)
	if err != nil {
		panic(err)
	}
	defer rocketmqProvider.ShutdownPushConsumer(consumer)

	p := provider{
		config:   cfg,
		tracing:  tracing,
		logger:   logger,
		pdb:      pdb,
		consumer: consumer,
		rocketmq: rocketmqProvider,
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

	// balance message
	balanceRepo := balance.NewRepository(p.pdb)
	balanceMessage := balance.NewMessage(p.consumer, balanceRepo, p.rocketmq, p.tracing, p.logger)
	if err := balanceMessage.Subscribe(); err != nil {
		panic(err)
	}
	if err := p.rocketmq.StartPushConsumer(p.consumer); err != nil {
		panic(err)
	}

	// customer
	customerRepo := customer.NewRepository(p.pdb)
	customerService := customer.NewService(customerRepo, p.logger)
	customer.RegisterHandlers(group, customerService)

	// Start server
	panic(e.Start(constants.RestPort))
}
