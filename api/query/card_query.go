package query

import (
	"strings"

	"gorm.io/gorm"
)

type CardQuery struct {
	Type     string  `form:"type"`
	Name     string  `form:"name"`
	MinPrice float32 `form:"minPrice,default=-1"`
	MaxPrice float32 `form:"maxPrice,default=-1"`
}

func (q *CardQuery) ApplyQueryF() func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		result := d.Where("card_type_id=? and LOWER(name) like ?", q.Type, "%"+strings.ToLower(q.Name)+"%")
		if q.MaxPrice != -1 {
			result = result.Where("price < ?", q.MaxPrice)
		}
		if q.MinPrice != -1 {
			result = result.Where("price > ?", q.MinPrice)
		}
		return result
	}
}
