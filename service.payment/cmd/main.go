package main

import (
	_ "context"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.payment/internal/payment"
	"os"
)

type provider struct {
	config   config.Config
	tracing  tracing.Provider
	logger   logger.Provider
	rocketmq rocketmq.Provider
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

	// init rocketmq
	rocketmq := rocketmq.NewProvider(cfg.RocketMQ)
	//if err = rocketMQProvider.CreateTopic(context.Background(), constants.RocketMQTopicPayment); err != nil {
	//	panic(err)
	//}

	p := provider{
		config:   cfg,
		tracing:  tracing,
		logger:   logger,
		rocketmq: rocketmq,
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

	payment.RegisterHandlers(group, p.rocketmq)

	// Start rest server
	panic(e.Start(constants.RestPort))
}
