package sku

type sku struct {
	ProductID string    `json:"productId"`
	Assets    skuAssets `json:"assets"`
	Attrs     []attr    `json:"attrs"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
}

type attr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type skuAssets struct {
	PPT []image `json:"ppt"`
}

type image struct {
	Src string `json:"src"`
}
