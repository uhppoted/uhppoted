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

type GetCardRequest struct {
	DeviceID   uint32
	CardNumber uint32
}

type GetCardResponse struct {
	MetaInfo interface{} `json:"meta-info,omitempty"`
	Device   Device      `json:"device"`
	Card     types.Card  `json:"card"`
}

func (m *GetCardResponse) SetMetaInfo(meta interface{}) {
	m.MetaInfo = meta
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
	MetaInfo interface{} `json:"meta-info,omitempty"`
	Device   Device      `json:"device"`
	Card     struct {
		CardNumber uint32 `json:"card-number"`
		Deleted    bool   `json:"deleted"`
	} `json:"card"`
}

func (m *DeleteCardResponse) SetMetaInfo(meta interface{}) {
	m.MetaInfo = meta
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

func (u *UHPPOTED) GetCard(ctx context.Context, request GetCardRequest) (*GetCardResponse, int, error) {
	u.debug(ctx, "get-card", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	cardID := request.CardNumber

	card, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetCardById(device, cardID)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error retrieving card %v from %v", cardID, device)
	}

	if card == nil {
		return nil, StatusNotFound, fmt.Errorf("No record for card %v on %v", cardID, device)
	}

	response := GetCardResponse{
		Device: Device{device},
		Card:   *card,
	}

	u.debug(ctx, "get-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
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
	cardNo := request.CardNumber

	deleted, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).DeleteCard(device, cardNo)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Error deleting card %v from %v", cardNo, device)
	}

	response := DeleteCardResponse{
		Device: Device{device},
		Card: struct {
			CardNumber uint32 `json:"card-number"`
			Deleted    bool   `json:"deleted"`
		}{
			CardNumber: cardNo,
			Deleted:    deleted.Succeeded,
		},
	}

	u.debug(ctx, "delete-card", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}
