package rest

import (
	"context"
	"errors"
	"github.com/uhppoted/uhppoted/src/uhppote"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net/http"
)

func getCards(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)

	N, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCards(deviceID)
	if err != nil {
		warn(ctx, deviceID, "get-cards", err)
		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}

	cards := make([]uint32, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByIndex(deviceID, index+1)
		if err != nil {
			warn(ctx, deviceID, "get-card-by-index", err)
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

func getCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardByID(deviceID, cardNumber)
	if err != nil {
		warn(ctx, deviceID, "get-card-by-id", err)
		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
		return
	}

	if card == nil {
		http.Error(w, "Card record does not exist", http.StatusNotFound)
		return
	}

	response := struct {
		Card types.Card `json:"card"`
	}{
		Card: *card,
	}

	reply(ctx, w, response)
}

func deleteCards(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCards(deviceID)
	if err != nil {
		warn(ctx, deviceID, "delete-cards", err)
		http.Error(w, "Error deleting cards", http.StatusInternalServerError)
		return
	}

	if !result.Succeeded {
		warn(ctx, deviceID, "delete-cards", errors.New("Request failed"))
		http.Error(w, "Error deleting cards", http.StatusInternalServerError)
		return
	}
}

func deleteCard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	cardNumber := ctx.Value("card-number").(uint32)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCard(deviceID, cardNumber)
	if err != nil {
		warn(ctx, deviceID, "delete-card-by-id", err)
		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
		return
	}

	if !result.Succeeded {
		warn(ctx, deviceID, "delete-card", errors.New("Request failed"))
		http.Error(w, "Error deleting card", http.StatusInternalServerError)
		return
	}
}
