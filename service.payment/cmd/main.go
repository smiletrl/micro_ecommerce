package main

import (
	"context"
	rocketmqLib "github.com/apache/rocketmq-client-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.payment/internal/payment"
	"os"
)

type provider struct {
	config   config.Config
	tracing  tracing.Provider
	logger   logger.Provider
	rocketmq rocketmqLib.Producer
	pdb      postgresql.Provider
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

	// init rocketmq
	rocketmqProvider := rocketmq.NewProvider(cfg.RocketMQ)
	producer, err := rocketmqProvider.CreateProducer(context.Background(), constants.RocketMQGroupPayment)
	if err != nil {
		panic(err)
	}
	defer rocketmqProvider.ShutdownProducer(producer)

	// init postgres
	pdb, err := postgresql.NewProvider(cfg, tracing)
	if err != nil {
		panic(err)
	}
	defer pdb.Close()

	p := provider{
		config:   cfg,
		tracing:  tracing,
		logger:   logger,
		rocketmq: producer,
		pdb:      pdb,
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

	paymentRepo := payment.NewRepository(p.pdb)
	paymentService := payment.NewService(paymentRepo, p.rocketmq, p.tracing, p.logger)
	payment.RegisterHandlers(group, paymentService)

	// Start rest server
	panic(e.Start(constants.RestPort))
}
