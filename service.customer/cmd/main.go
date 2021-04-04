package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/balance"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/customer"
)

func main() {
	// Echo instance
	e := echo.New()
	echoGroup := e.Group("api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	stage := os.Getenv(constants.Stage)
	config := config.Load(stage)
	db := dbcontext.InitDB(config, stage)

	// customer service
	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo)
	customer.RegisterHandlers(echoGroup, customerService)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
