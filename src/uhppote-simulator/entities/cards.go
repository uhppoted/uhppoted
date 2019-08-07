package entities

import (
	"uhppote/types"
)

type Card struct {
	CardNumber uint32         `json:"number"`
	From       types.Date     `json:"from"`
	To         types.Date     `json:"to"`
	Doors      map[uint8]bool `json:"doors"`
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

func (l *CardList) Delete(cardNumber uint32) bool {
	for ix, c := range *l {
		if c.CardNumber == cardNumber {
			copy((*l)[ix:], (*l)[ix+1:])
			(*l)[len(*l)-1] = nil
			*l = (*l)[:len(*l)-1]
			return true
		}
	}

	return false
}

func (l *CardList) DeleteAll() bool {
	old := *l
	*l = (*l)[:0]
	for ix, _ := range old {
		old[ix] = nil
	}

	return true
}
