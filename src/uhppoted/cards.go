package uhppoted

import (
	"context"
	"uhppote"
)

type CardList struct {
	Device struct {
		ID    uint32   `json:"id"`
		Cards []uint32 `json:"cards"`
	} `json:"device"`
}

func (u *UHPPOTED) GetCards(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-cards", rq)

	id, err := rq.DeviceId()
	if err != nil {
		u.warn(ctx, 0, "get-cards", err)
		u.oops(ctx, "get-cards", "Error retrieving cards (invalid device ID)", StatusBadRequest)
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

//func getCards(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	deviceId := ctx.Value("device-id").(uint32)
//
//	N, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCards(deviceId)
//	if err != nil {
//		warn(ctx, deviceId, "get-cards", err)
//		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
//		return
//	}
//
//	cards := make([]uint32, 0)
//
//	for index := uint32(0); index < N.Records; index++ {
//		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByIndex(deviceId, index+1)
//		if err != nil {
//			warn(ctx, deviceId, "get-card-by-index", err)
//			http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
//			return
//		}
//
//		cards = append(cards, record.CardNumber)
//	}
//
//	response := struct {
//		Cards []uint32 `json:"cards"`
//	}{
//		Cards: cards,
//	}
//
//	reply(ctx, w, response)
//}

//func getCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	deviceId := ctx.Value("device-id").(uint32)
//	cardNumber := ctx.Value("card-number").(uint32)
//
//	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardById(deviceId, cardNumber)
//	if err != nil {
//		warn(ctx, deviceId, "get-card-by-id", err)
//		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
//		return
//	}
//
//	if card == nil {
//		http.Error(w, "Card record does not exist", http.StatusNotFound)
//		return
//	}
//
//	response := struct {
//		Card types.Card `json:"card"`
//	}{
//		Card: *card,
//	}
//
//	reply(ctx, w, response)
//}
//
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
