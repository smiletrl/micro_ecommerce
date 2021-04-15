package cart

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
)

type ProductProxy interface {
	GetSKU(c echo.Context, skuID int64) (entity.SKU, error)
}
