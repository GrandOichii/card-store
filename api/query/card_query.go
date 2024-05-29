package query

type CardQuery struct {
	Raw string

	Type      string  `form:"type" url:"type"`
	Language  string  `form:"lang" url:"lang"`
	Name      string  `form:"name" url:"name"`
	Key       string  `form:"key" url:"key"`
	MinPrice  float32 `form:"minPrice,default=-1" url:"minPrice"`
	MaxPrice  float32 `form:"maxPrice,default=-1" url:"maxPrice"`
	Page      uint    `form:"page,default=1" url:"page"`
	Keywords  string  `form:"t" url:"keywords"`
	Expansion string  `form:"expansion" url:"expansion"`
}
