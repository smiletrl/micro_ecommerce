package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/smiletrl/micro_ecommerce/pkg/config"
	_ "github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/service.cart/internal/cart"
	productClient "github.com/smiletrl/micro_ecommerce/service.product/external/client"
	_ "os"
)

func main() {
	// Echo instance
	e := echo.New()
	echoGroup := e.Group("api/v1")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// initialize service
	/*
		stage := os.Getenv(constants.Stage)
		if stage == "" {
			stage = constants.StageLocal
		}
		config, err := config.Load(stage)
		if err != nil {
			panic(err)
		}*/
	/*db, err := dbcontext.InitDB(config)
	if err != nil {
		panic(err)
	}*/

	db := dbcontext.NewDBContext(nil)
	healthcheck.RegisterHandlers(e.Group(""), db)

	// product rpc client. Inject config
	// this client depends on the product rpc server. Maybe Build the client
	// on needed basis. When the client call is required, initialize the
	// client connection.
	// Or build the connection on the product package?
	// @todo test the connection. terminate the rpc server, and then reenable it, and see if it works
	// I assume this cart service needs restart too.
	// The connection should be built on product side
	pclient := productClient.NewClient()

	// cart
	cartRepo := cart.NewRepository(db)
	productProxy := product{pclient}
	cartService := cart.NewService(cartRepo, productProxy)
	cart.RegisterHandlers(echoGroup, cartService)

	// Start server
	e.Logger.Fatal(e.Start(":1325"))
}

// product proxy
type product struct {
	client productClient.Client
}

func (p product) GetDetail(c echo.Context, id int64) (entity.Product, error) {
	return p.client.GetProduct(id)
}
