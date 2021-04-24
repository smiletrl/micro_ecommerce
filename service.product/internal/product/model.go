package product

import "time"

// product models product
type product struct {
	Title    string        `json:"title"`
	Body     string        `json:"body"`
	Category string        `json:"category"`
	Assets   productAssets `json:"assets"`
	Variants variantConfig `json:"variants"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type productAssets struct {
	PPT       []image `json:"ppt"`
	Thumbnail image   `json:"thumbnail"`
}

type image struct {
	Src string `json:"src"`
}

type variantConfig struct {
	Attrs []attrConfig `json:"attrs"`
}

type attrConfig struct {
	Name string `json:"name"`
}
