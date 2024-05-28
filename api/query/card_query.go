package query

// keywords search example: mtg black alpha

type CardQuery struct {
	Raw       string
	Type      string  `form:"type"`
	Language  string  `form:"lang"`
	Name      string  `form:"name"`
	Key       string  `form:"key"`
	MinPrice  float32 `form:"minPrice,default=-1"`
	MaxPrice  float32 `form:"maxPrice,default=-1"`
	Page      uint    `form:"page,default=1"`
	Keywords  string  `form:"t"`
	Expansion string  `form:"expansion"`
}
