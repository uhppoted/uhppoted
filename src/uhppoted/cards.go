package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
	"uhppote/types"
)

type GetCardsResponse struct {
	Device struct {
		ID    uint32   `json:"id"`
		Cards []uint32 `json:"cards"`
	} `json:"device"`
}

type GetCardResponse struct {
	Device struct {
		ID   uint32     `json:"id"`
		Card types.Card `json:"card"`
	} `json:"device"`
}

type PutCardResponse struct {
	Device struct {
		ID         uint32 `json:"id"`
		CardNumber uint32 `json:"card-number"`
		Authorized bool   `json:"authorized"`
	} `json:"device"`
}

type DeleteCardRequest struct {
	DeviceID   uint32
	CardNumber uint32
}

type DeleteCardResponse struct {
	Device struct {
		ID         uint32 `json:"id"`
		CardNumber uint32 `json:"card-number"`
		Deleted    bool   `json:"deleted"`
	} `json:"device"`
}

type DeleteCardsResponse struct {
	Device struct {
		ID      uint32 `json:"id"`
		Deleted bool   `json:"deleted"`
	} `json:"device"`
}

func (u *UHPPOTED) GetCards(ctx context.Context, rq Request) {
	u.debug(ctx, "get-cards", rq)

	id, err := rq.DeviceID()
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

	response := GetCardsResponse{
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

func (u *UHPPOTED) DeleteCards(ctx context.Context, rq Request) {
	u.debug(ctx, "delete-cards", rq)

	id, err := rq.DeviceID()
	if err != nil {
		u.warn(ctx, 0, "delete-cards", err)
		u.oops(ctx, "delete-cards", "Missing/invalid device ID)", StatusBadRequest)
		return
	}

	deleted, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCards(*id)
	if err != nil {
		u.warn(ctx, *id, "delete-cards", err)
		u.oops(ctx, "delete-cards", "Error deleting cards", StatusInternalServerError)
		return
	}

	response := DeleteCardsResponse{
		struct {
			ID      uint32 `json:"id"`
			Deleted bool   `json:"deleted"`
		}{
			ID:      *id,
			Deleted: deleted.Succeeded,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) GetCard(ctx context.Context, rq Request) {
	u.debug(ctx, "get-card", rq)

	id, cardnumber, err := rq.DeviceCardID()
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

	response := GetCardResponse{
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

func (u *UHPPOTED) PutCard(ctx context.Context, rq Request) {
	u.debug(ctx, "put-card", rq)

	id, card, err := rq.DeviceCard()
	if err != nil {
		u.warn(ctx, 0, "put-card", err)
		u.oops(ctx, "put-card", "Missing/invalid device ID or card information)", StatusBadRequest)
		return
	}

	authorized, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).PutCard(*id, *card)
	if err != nil {
		u.warn(ctx, *id, "put-card", err)
		u.oops(ctx, "put-card", "Error adding/updating card", StatusInternalServerError)
		return
	}

	response := PutCardResponse{
		struct {
			ID         uint32 `json:"id"`
			CardNumber uint32 `json:"card-number"`
			Authorized bool   `json:"authorized"`
		}{
			ID:         *id,
			CardNumber: card.CardNumber,
			Authorized: authorized.Succeeded,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) DeleteCard(ctx context.Context, request DeleteCardRequest) (*DeleteCardResponse, int, error) {
	u.debug(ctx, "delete-card", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	card := request.CardNumber

	deleted, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCard(device, card)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error deleting card %v from %v", card, device)
	}

	response := DeleteCardResponse{
		struct {
			ID         uint32 `json:"id"`
			CardNumber uint32 `json:"card-number"`
			Deleted    bool   `json:"deleted"`
		}{
			ID:         device,
			CardNumber: card,
			Deleted:    deleted.Succeeded,
		},
	}

	u.debug(ctx, "delete-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}
