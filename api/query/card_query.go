package query

import (
	"strings"

	"gorm.io/gorm"
)

type CardQuery struct {
	Type string `form:"type"`
	Name string `form:"name"`
}

func (q *CardQuery) ApplyQueryF() func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("card_type_id=? and LOWER(name) like ?", q.Type, "%"+strings.ToLower(q.Name)+"%")
	}
}
