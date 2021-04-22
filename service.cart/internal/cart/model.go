package cart

import (
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
)

type cartItem struct {
	Quantity int  `json:"quantity"`
	Valid    bool `json:"valid"`
	entity.SkuProperty
}
