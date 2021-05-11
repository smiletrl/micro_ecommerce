package cart

import (
	"context"
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
)

type ProductProxy interface {
	GetSkuStock(c context.Context, skuID string) (int, error)
	GetSkuProperties(c context.Context, skuIDs []string) ([]entity.SkuProperty, error)
}
