package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"github.com/smiletrl/micro_ecommerce/pkg/kafka"
	loggerd "github.com/smiletrl/micro_ecommerce/service.logger/internal/logger"
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

	// initialize service
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	config, err := config.Load(stage)
	if err != nil {
		panic(err)
	}

	healthcheck.RegisterHandlers(e.Group(""))

	// kafka message
	kafkaProvider := kafka.NewProvider(config.Kafka, sugar)

	partition := 0
	topic := constants.KafkaTopic("logger")
	err = kafkaProvider.CreateTopic(context.Background(), topic, partition)
	if err != nil {
		panic(err)
	}
	jwtService := jwt.NewService(config.JwtSecret)

	// kafka message
	err = loggerd.Consume(config.Kafka, sugar, topic, partition)
	if err != nil {
		panic(err)
	}
	loggerRepo := loggerd.NewRepository(kafkaProvider)
	loggerService := loggerd.NewService(loggerRepo, sugar)
	loggerd.RegisterHandlers(echoGroup, loggerService, jwtService)

	// Start server
	e.Logger.Fatal(e.Start(constants.RestPort))
}
