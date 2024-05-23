package dto

import (
	"store.api/model"
	"store.api/utility"
)

type GetCollection struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Cards       []*GetCardSlot `json:"cards"`
}

func NewGetCollection(col *model.Collection) *GetCollection {
	return &GetCollection{
		ID:          col.ID,
		Name:        col.Name,
		Description: col.Description,
		Cards: utility.MapSlice(
			col.Cards,
			func(c model.CardSlot) *GetCardSlot {
				return NewGetCardSlot(&c)
			},
		),
	}
}
