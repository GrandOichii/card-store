package query

import (
	"strings"

	"gorm.io/gorm"
)

type CardQuery struct {
	Type     string  `form:"type"`
	Language string  `form:"lang"`
	Name     string  `form:"name"`
	MinPrice float32 `form:"minPrice,default=-1"`
	MaxPrice float32 `form:"maxPrice,default=-1"`
}

func (q *CardQuery) ApplyQueryF() func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		result := d.Where("LOWER(name) like ?", "%"+strings.ToLower(q.Name)+"%")
		if len(q.Type) > 0 {
			result = result.Where("card_type_id=?", q.Type)
		}
		if len(q.Language) > 0 {
			result = result.Where("language_id=?", q.Language)
		}
		if q.MaxPrice != -1 {
			result = result.Where("price < ?", q.MaxPrice)
		}
		if q.MinPrice != -1 {
			result = result.Where("price > ?", q.MinPrice)
		}
		return result
	}
}
