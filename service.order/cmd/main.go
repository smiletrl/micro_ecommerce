package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/service.order/internal/order"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()
	//echoGroup := e.Group("/api/v1")

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
	db, err := dbcontext.InitDB(config)
	if err != nil {
		panic(err)
	}

	healthcheck.RegisterHandlers(e.Group(""), db)

	// rocketMQ message
	err = order.Consume(config.RocketMQ)
	if err != nil {
		panic(err)
	}
	// Start server
	e.Logger.Fatal(e.Start(constants.RestPort))
}
