package rest

import (
	"context"
	"net/http"
	"uhppote"
)

func getCards(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)

	N, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCards(deviceId)
	if err != nil {
		warn(ctx, deviceId, "get-cards", err)
		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}

	cards := make([]uint32, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByIndex(deviceId, index+1)
		if err != nil {
			warn(ctx, deviceId, "get-card-by-index", err)
			http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
			return
		}

		cards = append(cards, record.CardNumber)
	}

	response := struct {
		Cards []uint32 `json:"cards"`
	}{
		Cards: cards,
	}

	reply(ctx, w, response)
}
