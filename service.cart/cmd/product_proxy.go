package main

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
	productClient "github.com/smiletrl/micro_ecommerce/service.product/external/client"
)

// product proxy
type product struct {
	client productClient.Client
}

func (p product) GetDetail(c echo.Context, id int64) (entity.Product, error) {
	return p.client.GetProduct(id)
}
