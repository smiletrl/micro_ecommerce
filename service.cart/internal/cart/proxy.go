package cart

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
)

type ProductProxy interface {
	GetSkuStock(c echo.Context, skuID string) (int, error)
	GetSkuProperties(c echo.Context, skuIDs []string) ([]entity.SkuProperty, error)
}
