package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/service.payment/internal/payment"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()
	echoGroup := e.Group("/api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialize service
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	config, err := config.Load(stage)
	if err != nil {
		panic(err)
	}
	db, err := postgresql.InitDB(config)
	if err != nil {
		panic(err)
	}

	healthcheck.RegisterHandlers(e.Group(""))

	paymentRepo := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepo)

	rocketMQProvider := rocketmq.NewProvider(config.RocketMQ)
	if err = rocketMQProvider.CreateTopic(context.Background(), constants.RocketMQTopicPayment); err != nil {
		panic(err)
	}
	payment.RegisterHandlers(echoGroup, paymentService, rocketMQProvider)

	// Start rest server
	e.Logger.Fatal(e.Start(constants.RestPort))
}
