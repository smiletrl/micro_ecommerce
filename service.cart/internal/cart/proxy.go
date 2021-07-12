package cart

import (
	"context"

	"github.com/smiletrl/micro_ecommerce/pkg/entity"
)

type ProductProxy interface {
	GetSkuStock(c context.Context, skuID string) (int, error)
	GetSkuProperties(c context.Context, skuIDs []string) ([]entity.SkuProperty, error)
}

type mockProduct struct{}

func newMockProduct() ProductProxy {
	return mockProduct{}
}

func (m mockProduct) GetSkuStock(c context.Context, skuID string) (int, error) {
	return 12, nil
}

func (m mockProduct) GetSkuProperties(c context.Context, skuIDs []string) ([]entity.SkuProperty, error) {
	return []entity.SkuProperty{
		entity.SkuProperty{
			SkuID:      "12",
			Title:      "mac",
			Price:      1200,
			Attributes: "xl",
			Thumbnail:  "xx.ng",
			Stock:      12,
		},
	}, nil
}
