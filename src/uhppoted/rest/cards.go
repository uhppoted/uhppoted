package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"uhppote"
)

func getCards(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)
	log := ctx.Value("log").(*log.Logger)

	N, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCards(deviceId)
	if err != nil {
		log.Printf("WARN: [getCards] %v\n", err)
		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}

	cards := make([]uint32, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByIndex(deviceId, index+1)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving cards: %v", err), http.StatusInternalServerError)
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
