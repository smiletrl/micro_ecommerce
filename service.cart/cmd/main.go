package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/service.cart/internal/cart"
	productClient "github.com/smiletrl/micro_ecommerce/service.product/external"
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
	db, err := dbcontext.InitDB(config)
	if err != nil {
		panic(err)
	}

	rocketmq.Start()

	healthcheck.RegisterHandlers(e.Group(""), db)

	// Product rpc client. Inject config
	pclient := productClient.NewClient()

	// cart
	cartRepo := cart.NewRepository(db)
	productProxy := product{pclient}
	cartService := cart.NewService(cartRepo, productProxy)
	cart.RegisterHandlers(echoGroup, cartService)

	// Start server
	e.Logger.Fatal(e.Start(constants.RestPort))
}

// product proxy
type product struct {
	client productClient.Client
}

func (p product) GetSKU(c echo.Context, skuID int64) (entity.SKU, error) {
	return p.client.GetSKU(skuID)
}
