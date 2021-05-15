package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"github.com/smiletrl/micro_ecommerce/pkg/kafka"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	loggerd "github.com/smiletrl/micro_ecommerce/service.logger/internal/logger"
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

	// init tracing
	tracingProvider := tracing.NewProvider()
	closer, err := tracingProvider.SetupTracer(constants.TracingPayment, cfg)
	if err != nil {
		logger.Fatal("tracing", err)
	}
	defer closer.Close()

	// echo instance
	e := echo.New()
	echoGroup := e.Group("/api/v1")

	// middleware
	e.Use(accesslog.Middleware(logger))
	e.Use(tracingProvider.Middleware(logger))
	e.Use(errors.Middleware(logger))

	healthcheck.RegisterHandlers(e.Group(""))

	// kafka message
	kafkaProvider := kafka.NewProvider(cfg.Kafka, logger)

	partition := 0
	topic := constants.KafkaTopic("logger")
	err = kafkaProvider.CreateTopic(context.Background(), topic, partition)
	if err != nil {
		panic(err)
	}
	jwtProvider := jwt.NewProvider(cfg.JwtSecret)

	// kafka message
	err = loggerd.Consume(cfg.Kafka, logger, topic, partition)
	if err != nil {
		panic(err)
	}
	loggerRepo := loggerd.NewRepository(kafkaProvider)
	loggerService := loggerd.NewService(loggerRepo, logger)
	loggerd.RegisterHandlers(echoGroup, loggerService, jwtProvider)

	// Start server
	panic(e.Start(constants.RestPort))
}
