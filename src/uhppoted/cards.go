package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
	"uhppote/types"
)

type CardList struct {
	Device struct {
		ID    uint32   `json:"id"`
		Cards []uint32 `json:"cards"`
	} `json:"device"`
}

type Card struct {
	Device struct {
		ID   uint32     `json:"id"`
		Card types.Card `json:"card"`
	} `json:"device"`
}

func (u *UHPPOTED) GetCards(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-cards", rq)

	id, err := rq.DeviceId()
	if err != nil {
		u.warn(ctx, 0, "get-cards", err)
		u.oops(ctx, "get-cards", "Missing/invalid device ID)", StatusBadRequest)
		return
	}

	N, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCards(*id)
	if err != nil {
		u.warn(ctx, *id, "get-cards", err)
		u.oops(ctx, "get-cards", "Error retrieving cards", StatusInternalServerError)
		return
	}

	cards := make([]uint32, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByIndex(*id, index+1)
		if err != nil {
			u.warn(ctx, *id, "get-cards", err)
			u.oops(ctx, "get-cards", "Error retrieving cards", StatusInternalServerError)
			return
		}

		cards = append(cards, record.CardNumber)
	}

	response := CardList{
		struct {
			ID    uint32   `json:"id"`
			Cards []uint32 `json:"cards"`
		}{
			ID:    *id,
			Cards: cards,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) GetCard(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-card", rq)

	id, cardnumber, err := rq.DeviceCard()
	if err != nil {
		u.warn(ctx, 0, "get-card", err)
		u.oops(ctx, "get-card", "Missing/invalid device ID or card number)", StatusBadRequest)
		return
	}

	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardById(*id, *cardnumber)
	if err != nil {
		u.warn(ctx, *id, "get-card", err)
		u.oops(ctx, "get-card", "Error retrieving card", StatusInternalServerError)
		return
	}

	if card == nil {
		u.warn(ctx, *id, "get-card", fmt.Errorf("No record for card %d", *cardnumber))
		u.oops(ctx, "get-card", fmt.Sprintf("No record for card %d", *cardnumber), StatusNotFound)
		return
	}

	response := Card{
		struct {
			ID   uint32     `json:"id"`
			Card types.Card `json:"card"`
		}{
			ID:   *id,
			Card: *card,
		},
	}

	u.reply(ctx, response)
}

//func deleteCards(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	deviceId := ctx.Value("device-id").(uint32)
//
//	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCards(deviceId)
//	if err != nil {
//		warn(ctx, deviceId, "delete-cards", err)
//		http.Error(w, "Error deleting cards", http.StatusInternalServerError)
//		return
//	}
//
//	if !result.Succeeded {
//		warn(ctx, deviceId, "delete-cards", errors.New("Request failed"))
//		http.Error(w, "Error deleting cards", http.StatusInternalServerError)
//		return
//	}
//}
//
//func deleteCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	deviceId := ctx.Value("device-id").(uint32)
//	cardNumber := ctx.Value("card-number").(uint32)
//
//	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCard(deviceId, cardNumber)
//	if err != nil {
//		warn(ctx, deviceId, "delete-card-by-id", err)
//		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
//		return
//	}
//
//	if !result.Succeeded {
//		warn(ctx, deviceId, "delete-card", errors.New("Request failed"))
//		http.Error(w, "Error deleting card", http.StatusInternalServerError)
//		return
//	}
//}
