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
	closer, err := tracingProvider.SetupTracer("cart", cfg)
	if err != nil {
		logger.Fatal("tracing", err)
	}
	defer closer.Close()

	// middleware
	e.Use(accesslog.Middleware(logger))
	e.Use(tracingProvider.Middleware(logger))
	e.Use(errors.Middleware(logger))

	healthcheck.RegisterHandlers(e.Group(""))

	rocketMQProvider := rocketmq.NewProvider(cfg.RocketMQ)
	//if err = rocketMQProvider.CreateTopic(context.Background(), constants.RocketMQTopicPayment); err != nil {
	//	panic(err)
	//}
	payment.RegisterHandlers(echoGroup, rocketMQProvider)

	// Start rest server
	e.Logger.Fatal(e.Start(constants.RestPort))
}
