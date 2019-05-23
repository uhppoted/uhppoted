package entities

import (
	"uhppote/types"
)

type Card struct {
	CardNumber uint32     `json:"number"`
	From       types.Date `json:"from"`
	To         types.Date `json:"to"`
	Door1      bool       `json:"door-1"`
	Door2      bool       `json:"door-2"`
	Door3      bool       `json:"door-3"`
	Door4      bool       `json:"door-4"`
}

type CardList []*Card

// TODO: implement Marshal/Unmarshal
func (l *CardList) Put(card *Card) {
	if card != nil {
		index := -1
		for i, c := range *l {
			if c.CardNumber == card.CardNumber {
				index = i
				break
			}
		}

		if index == -1 {
			*l = append(*l, card)
		} else {
			(*l)[index] = card
		}
	}
}
