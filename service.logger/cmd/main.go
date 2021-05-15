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

type provider struct {
	config  config.Config
	logger  logger.Provider
	tracing tracing.Provider
	kafka   kafka.Provider
	jwt     jwt.Provider
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

	p := provider{
		config:  cfg,
		tracing: tracing,
		logger:  logger,
		kafka:   kafkaProvider,
		jwt:     jwtProvider,
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

	// logger
	loggerRepo := loggerd.NewRepository(p.kafka)
	loggerService := loggerd.NewService(loggerRepo, p.logger)
	loggerd.RegisterHandlers(group, loggerService, p.jwt)

	// Start rest server
	panic(e.Start(constants.RestPort))
}
