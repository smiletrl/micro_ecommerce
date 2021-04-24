package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	_ "github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
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
	healthcheck.RegisterHandlers(e.Group(""))

	// redis
	redisClient := redis.New(config)

	jwtService := jwt.NewService(config.JwtSecret)

	// Product rpc client. Inject config
	pclient := productClient.NewClient()

	// cart
	cartRepo := cart.NewRepository(redisClient)

	productProxy := product{pclient}
	cartService := cart.NewService(cartRepo, productProxy)
	cart.RegisterHandlers(echoGroup, cartService, jwtService)

	// Start server
	e.Logger.Fatal(e.Start(constants.RestPort))
}

// product proxy
type product struct {
	client productClient.Client
}

func (p product) GetSkuStock(c echo.Context, skuID string) (int, error) {
	// maybe we want to add timeout for this request in case this request just hangs on.
	return p.client.GetSkuStock(c, skuID)
}

func (p product) GetSkuProperties(c echo.Context, skuIDs []string) ([]entity.SkuProperty, error) {
	return p.client.GetSkuProperties(c, skuIDs)
}
