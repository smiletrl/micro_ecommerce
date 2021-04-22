package cart

import (
	"github.com/smiletrl/micro_ecommerce/pkg/entity"
)

type cartItem struct {
	Quantity int  `json:"quantity"`
	Valid    bool `json:"valid"`
	entity.SkuProperty
}

/*
type cartItemDetail struct {
	cartItem
	Title      string `json:"title"`
	Price      int    `json:"price"`
	Attributes string `json:"attributes"`
	Thumbnail  string `json:"thumbnail"`
	Valid      bool   `json:"valid"`
}

type cart []cartItemDetail
*/
